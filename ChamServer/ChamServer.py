import socketserver,binascii,base64,socket,logging,redis,google.protobuf
from ChamDataFromat.data_pb2 import *
import MAPCH, json

class Cache():
    def __init__(self, host='localhost', port=6379):
        logging.info("connect to cache database {}:{}".format(host, port))
        self.server = redis.Redis(host=host, port=port, db=0)
        logging.info("connect to cache database successful")

    '''
    set key:value pair to database.
    return true/false
    '''
    def set(self, key, value)->bool:
        return self.server.set(key, value)

    '''
    get value of key in database.
    return (false,None) means there is no that key in database.
    otherwise return (true, value).
    '''
    def get(self, key)->(bool, object):
        value = self.server.get(key)
        if None == value:
            return False, None
        else:
            return True, value

class MyTCPHandler(socketserver.BaseRequestHandler):
    """
    The request handler class for our server.

    It is instantiated once per connection to the server, and must
    override the handle() method to implement communication to the
    client.
    """
    def handle(self):
        # self.request is the TCP socket connected to the client

        logging.info("recv from {}:{}".format(self.client_address[0],self.client_address[1]))
        m = self.request.recv(1024*10)
        tf = Transfer()
        tf.ParseFromString(m)

        if tf.ttype == TransferType.Hash:
            logging.info("handle chamhash function")
            recvdata = self.handleHashFunction(tf.tdata)
            recvtf = Transfer(ttype=TransferType.HashRecv, tdata=recvdata)
            self.request.sendall(recvtf.SerializeToString())
        elif tf.ttype == TransferType.Check:
            logging.info("handle chamHash check")
            recvdata = self.handleCheckFunction(tf.tdata)
            recvtf = Transfer(ttype=TransferType.CheckRecv, tdata=recvdata)
            self.request.sendall(recvtf.SerializeToString())
        elif tf.ttype == TransferType.Adapt:
            logging.info("handle ChamHash Adapt")
            recvdata = self.handleAdaptFunction(tf.tdata)
            recvtf = Transfer(ttype=TransferType.AdaptRecv, tdata=recvdata)
            self.request.sendall(recvtf.SerializeToString())
        else :
            logging.info("not such that function")
        # judge message in redis database or not
        # self.request.close()

    def handleCheckFunction(self, mhbytes):
        logging.info("handle check function")
        mh = CheckSet()
        mh.ParseFromString(mhbytes)

        state = MAPCH.check(mh.m.mes, json.loads(mh.ch.helperData))
        ck = CheckState(check = state)
        value = ck.SerializeToString()
        logging.info("handle message check state: {}".format(str(state)))
        return value

    def handleAdaptFunction(self, mmhbytes):
        logging.info("handle adapt function")
        mmh = AdaptSet()
        mmh.ParseFromString(mmhbytes)

        h = MAPCH.collision(mmh.m1.mes, mmh.m2.mes, json.loads(mmh.ch.helperData))

        ch = ChamHash(hash = h['h'], helperData = json.dumps(h))
        value = ch.SerializeToString()
        return value

    def handleHashFunction(self, mbytes):

        logging.debug("handle hash function with message {}".format(mbytes))
        m = Message()
        m.ParseFromString(mbytes)

        global Cache_Database
        b, value = Cache_Database.get(m.mes)
        if b == True:
            logging.info("cache hit in message: {}".format(m.mes))
            return value
        else:
            logging.info("fail to hit message in cache")

        h = MAPCH.hash(m.mes)
        ch = ChamHash(hash = h['h'], helperData = json.dumps(h))

        value = ch.SerializeToString()
        if not Cache_Database.set(m.mes, value):
            logging.error("redis database store fail")
        return value

if __name__ == "__main__":
    HOST, PORT = "0.0.0.0", 1234

    # config log level
    logging.basicConfig(level=logging.INFO)
    Cache_Database = Cache()
    socketserver.TCPServer.allow_reuse_address = True
    #socketserver.TCPServer.timeout = 1
    server = socketserver.TCPServer((HOST, PORT), MyTCPHandler)
    # Activate the server; this will keep running until you
    # interrupt the program with Ctrl-C
    server.serve_forever()

