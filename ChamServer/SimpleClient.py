import socket
from ChamDataFromat.data_pb2 import *

def Hash(m):
    s = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
    s.connect(("localhost",1234))

    tf = Transfer(ttype=TransferType.Hash, tdata=m.SerializeToString())
    s.send(tf.SerializeToString())
    data = s.recv(1024*10)
    s.close()

    tf = Transfer()
    tf.ParseFromString(data)
    ch = ChamHash()
    ch.ParseFromString(tf.tdata)
    return ch

def Check(m, ch):
    s = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
    checkdata = CheckSet(m = m, ch = ch).SerializeToString()

    checktf = Transfer(ttype = TransferType.Check, tdata = checkdata)
    s.connect(("localhost",1234))
    s.send(checktf.SerializeToString())
    data = s.recv(1024*10)
    s.close()

    tf = Transfer()
    tf.ParseFromString(data)
    cs  = CheckState()
    cs.ParseFromString(tf.tdata)

    return cs.check

def Adapt(m1, m2, h):
    adaptdata = AdaptSet(m1 = m1, m2=m2, ch=h).SerializeToString()
    adapttf = Transfer(ttype = TransferType.Adapt, tdata = adaptdata)

    s = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
    s.connect(("localhost",1234))
    s.send(adapttf.SerializeToString())
    data = s.recv(1024*10)
    s.close()

    tf = Transfer()
    tf.ParseFromString(data)
    newch  = ChamHash()
    newch.ParseFromString(tf.tdata)

    return newch

if __name__=="__main__":
    m1 = Message(mes = b'123')
    m2 = Message(mes = b'234')
    ch = Hash(m1)

    if Check(m1, ch):
        print("check PASS")
    else:
        print("check FAIL")

    newch = Adapt(m1, m2, ch)

    if Check(m2, newch):
        print("adapt PASS")
    else:
        print("adapt FAIL")

