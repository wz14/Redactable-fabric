import chamwithemp
import MAABE
import re,json
import pickle
from charm.toolbox.integergroup import integer
from charm.toolbox.pairinggroup import PairingGroup,GT
from charm.toolbox.symcrypto import AuthenticatedCryptoAbstraction,SymmetricCryptoAbstraction
from charm.core.math.pairing import hashPair as extractor

groupObj = PairingGroup('SS512')

maabe = MAABE.MaabeRW15(groupObj)
chamHash = chamwithemp.Chamwithemp()
groupInt = chamwithemp.group

def merge_dicts(*dict_args):
    """
    Given any number of dicts, shallow copy and merge into a new dict,
    precedence goes to key value pairs in latter dicts.
    """
    result = {}
    for dictionary in dict_args:
        result.update(dictionary)
    return result

# get pp
pp_file = open("pp_output.txt","rb")
pp_unpi = pickle.load(pp_file)
g1_deseri = groupObj.deserialize(pp_unpi['g1'])
g2_deseri = groupObj.deserialize(pp_unpi['g2'])
egg_deseri = groupObj.deserialize(pp_unpi['egg'])
public_parameters = {'g1':g1_deseri,'g2':g2_deseri,'egg':egg_deseri}
#print(public_parameters)

# get maabepk
maabepk_file = open("maabepk_output.txt","rb")
maabepk_unpi = pickle.load(maabepk_file)
egga1_deseri = groupObj.deserialize(maabepk_unpi['UT']['egga'])
gy1_deseri = groupObj.deserialize(maabepk_unpi['UT']['gy'])
egga2_deseri = groupObj.deserialize(maabepk_unpi['OU']['egga'])
gy2_deseri = groupObj.deserialize(maabepk_unpi['OU']['gy'])
pk1 = {'egga':egga1_deseri,'gy':gy1_deseri}
pk2 = {'egga':egga2_deseri,'gy':gy2_deseri}
maabepk = {'UT':pk1, 'OU':pk2}

# get maabesk
maabesk_file = open("maabesk_output.txt","rb")
maabesk_unpi = pickle.load(maabesk_file)
alpha1_deseri = groupObj.deserialize(maabesk_unpi['UT']['alpha'])
y1_deseri = groupObj.deserialize(maabesk_unpi['UT']['y'])
alpha2_deseri = groupObj.deserialize(maabesk_unpi['OU']['alpha'])
y2_deseri = groupObj.deserialize(maabesk_unpi['OU']['y'])
sk1 = {'alpha':alpha1_deseri,'y':y1_deseri}
sk2 = {'alpha':alpha2_deseri,'y':y2_deseri}
maabesk = {'UT':sk1,'OU':sk2}

#chamhash key init
#(pk, sk) = chamHash.keygen(1024)
pk = {'secparam': 1024, 'N': integer(20639501445490619499571473231962238639613850200609723129865274973389207169386019790081317166801377585852101611924513018118076503220272505269101123185483478827902539989267940452253498331363934637463307815140081875402144529850087490107865496237873649044005216754760322206792335556380079175312884873545407803758528692484712100168089812574325930202675592634522400631602766263795012899595172109896928945224406833545127918203360663537284395229838935471620482704282064477409110985233629664210200323159926539979025200450130171953285108024706981095404167217773267631722062136111440636883911386048815309175875991313022494471621),
'phi_N': integer(20639501445490619499571473231962238639613850200609723129865274973389207169386019790081317166801377585852101611924513018118076503220272505269101123185483478827902539989267940452253498331363934637463307815140081875402144529850087490107865496237873649044005216754760322206792335556380079175312884873545407803758240912083877065263368139473545769655301188586129214282684729097512547789322898451498171540561242228826274490229393652275591679714621114555989685997248623325129278710182201931402727164911921144228799194648263502895573605624620716960707464670951786841583960602580888653518660736530954570818782522777783140038144)}
sk = {'p': integer(135834871082057596556156087102880695603597521170179782120553666786815062117471253900663075002935794606928067261571608422525654781342120345661492724826043241007133331185364909682061795645464088333161857722263310721335141129059734282105483228532633124181853501522026119051116547291556472645868866042715094037889),
'q': integer(151945529752977308165517013677279851770806527223006566797483499495650048154802404498094329660228810111925360712395402839167060733875700569969303982207397911272698943866062823125411362602541307417064148079603358336376361271026529852591219318288847665956248032008525864314134102226304265711224602492524260395589)}

# user's gid and attribute, access policy definitions
gid = "bob"
user_attr1 = ['STUDENT@UT']
user_attr2 = ['STUDENT@OU']
access_policy = '((STUDENT@UT or PROFESSOR@OU) and (STUDENT@UT or MASTERS@OU))'

# get user's attribute keys
usersk_file = open("usersk_output.txt","rb")
usersk_unpi = pickle.load(usersk_file)
K1_deseri = groupObj.deserialize(usersk_unpi['keys']['STUDENT@UT']['K'])
KP1_deseri = groupObj.deserialize(usersk_unpi['keys']['STUDENT@UT']['KP'])
usersk1_deseri = {'STUDENT@UT':{'K':K1_deseri,'KP':KP1_deseri}}

K2_deseri = groupObj.deserialize(usersk_unpi['keys']['STUDENT@OU']['K'])
KP2_deseri = groupObj.deserialize(usersk_unpi['keys']['STUDENT@OU']['KP'])
usersk2_deseri = {'STUDENT@OU':{'K':K2_deseri,'KP':KP2_deseri}}
user_sk = {'GID':'bob', 'keys':merge_dicts(usersk1_deseri, usersk2_deseri)}

def cut_text(text,lenth): 
    textArr = re.findall('.{'+str(lenth)+'}', text) 
    textArr.append(text[(len(textArr)*lenth):]) 
    return textArr


def serializehash(h):
    # return bytes type h
    C0_seri =  groupObj.serialize(h['cipher']['rkc']['C0']).decode()
    C11_seri = groupObj.serialize(h['cipher']['rkc']['C1']['STUDENT@UT_0']).decode()
    C12_seri = groupObj.serialize(h['cipher']['rkc']['C1']['PROFESSOR@OU']).decode()
    C13_seri = groupObj.serialize(h['cipher']['rkc']['C1']['STUDENT@UT_1']).decode()
    C14_seri = groupObj.serialize(h['cipher']['rkc']['C1']['MASTERS@OU']).decode()
    C1_seri = {'STUDENT@UT_0':C11_seri, 'PROFESSOR@OU':C12_seri,'STUDENT@UT_1':C13_seri,'MASTERS@OU':C14_seri}

    C21_seri = groupObj.serialize(h['cipher']['rkc']['C2']['STUDENT@UT_0']).decode()
    C22_seri = groupObj.serialize(h['cipher']['rkc']['C2']['PROFESSOR@OU']).decode()
    C23_seri = groupObj.serialize(h['cipher']['rkc']['C2']['STUDENT@UT_1']).decode()
    C24_seri = groupObj.serialize(h['cipher']['rkc']['C2']['MASTERS@OU']).decode()
    C2_seri = {'STUDENT@UT_0':C21_seri, 'PROFESSOR@OU':C22_seri,'STUDENT@UT_1':C23_seri,'MASTERS@OU':C24_seri}


    C31_seri = groupObj.serialize(h['cipher']['rkc']['C3']['STUDENT@UT_0']).decode()
    C32_seri = groupObj.serialize(h['cipher']['rkc']['C3']['PROFESSOR@OU']).decode()
    C33_seri = groupObj.serialize(h['cipher']['rkc']['C3']['STUDENT@UT_1']).decode()
    C34_seri = groupObj.serialize(h['cipher']['rkc']['C3']['MASTERS@OU']).decode()
    C3_seri = {'STUDENT@UT_0':C31_seri, 'PROFESSOR@OU':C32_seri,'STUDENT@UT_1':C33_seri,'MASTERS@OU':C34_seri}

    C41_seri = groupObj.serialize(h['cipher']['rkc']['C4']['STUDENT@UT_0']).decode()
    C42_seri = groupObj.serialize(h['cipher']['rkc']['C4']['PROFESSOR@OU']).decode()
    C43_seri = groupObj.serialize(h['cipher']['rkc']['C4']['STUDENT@UT_1']).decode()
    C44_seri = groupObj.serialize(h['cipher']['rkc']['C4']['MASTERS@OU']).decode()
    C4_seri = {'STUDENT@UT_0':C41_seri, 'PROFESSOR@OU':C42_seri,'STUDENT@UT_1':C43_seri,'MASTERS@OU':C44_seri}

    maabect_seri = {'policy':h['cipher']['rkc']['policy'], 'C0':C0_seri, 'C1':C1_seri, 'C2':C2_seri, 'C3':C3_seri, 'C4':C4_seri}

    ct_seri = {'rkc':maabect_seri,'ec':h['cipher']['ec']}

    h_seri = groupInt.serialize(h['h']).decode()
    r_seri = groupInt.serialize(h['r']).decode()
    N1_seri = groupInt.serialize(h['N1']).decode()
    e_seri = groupInt.serialize(h['e']).decode()
    hash_seri = {'h': h_seri, 'r': r_seri, 'cipher':ct_seri, 'N1': N1_seri, 'e': e_seri}
    return hash_seri

def deserializehash(h):
    # h is bytes type, so turn it into the origin type (integer or pairing)
    h_deseri = groupInt.deserialize(h['h'].encode())
    r_deseri = groupInt.deserialize(h['r'].encode())
    N1_deseri = groupInt.deserialize(h['N1'].encode())
    e_deseri = groupInt.deserialize(h['e'].encode())

    C0_deseri = groupObj.deserialize(h['cipher']['rkc']['C0'].encode())
    C11_deseri = groupObj.deserialize(h['cipher']['rkc']['C1']['STUDENT@UT_0'].encode())
    C12_deseri = groupObj.deserialize(h['cipher']['rkc']['C1']['PROFESSOR@OU'].encode())
    C13_deseri = groupObj.deserialize(h['cipher']['rkc']['C1']['STUDENT@UT_1'].encode())
    C14_deseri = groupObj.deserialize(h['cipher']['rkc']['C1']['MASTERS@OU'].encode())
    C1_deseri = {'STUDENT@UT_0':C11_deseri, 'PROFESSOR@OU':C12_deseri,'STUDENT@UT_1':C13_deseri,'MASTERS@OU':C14_deseri}

    C21_deseri = groupObj.deserialize(h['cipher']['rkc']['C2']['STUDENT@UT_0'].encode())
    C22_deseri = groupObj.deserialize(h['cipher']['rkc']['C2']['PROFESSOR@OU'].encode())
    C23_deseri = groupObj.deserialize(h['cipher']['rkc']['C2']['STUDENT@UT_1'].encode())
    C24_deseri = groupObj.deserialize(h['cipher']['rkc']['C2']['MASTERS@OU'].encode())
    C2_deseri = {'STUDENT@UT_0':C21_deseri, 'PROFESSOR@OU':C22_deseri,'STUDENT@UT_1':C23_deseri,'MASTERS@OU':C24_deseri}


    C31_deseri = groupObj.deserialize(h['cipher']['rkc']['C3']['STUDENT@UT_0'].encode())
    C32_deseri = groupObj.deserialize(h['cipher']['rkc']['C3']['PROFESSOR@OU'].encode())
    C33_deseri = groupObj.deserialize(h['cipher']['rkc']['C3']['STUDENT@UT_1'].encode())
    C34_deseri = groupObj.deserialize(h['cipher']['rkc']['C3']['MASTERS@OU'].encode())
    C3_deseri = {'STUDENT@UT_0':C31_deseri, 'PROFESSOR@OU':C32_deseri,'STUDENT@UT_1':C33_deseri,'MASTERS@OU':C34_deseri}

    C41_deseri = groupObj.deserialize(h['cipher']['rkc']['C4']['STUDENT@UT_0'].encode())
    C42_deseri = groupObj.deserialize(h['cipher']['rkc']['C4']['PROFESSOR@OU'].encode())
    C43_deseri = groupObj.deserialize(h['cipher']['rkc']['C4']['STUDENT@UT_1'].encode())
    C44_deseri = groupObj.deserialize(h['cipher']['rkc']['C4']['MASTERS@OU'].encode())
    C4_deseri = {'STUDENT@UT_0':C41_deseri, 'PROFESSOR@OU':C42_deseri,'STUDENT@UT_1':C43_deseri,'MASTERS@OU':C44_deseri}

    maabect_deseri = {'policy':h['cipher']['rkc']['policy'], 'C0':C0_deseri, 'C1':C1_deseri, 'C2':C2_deseri, 'C3':C3_deseri, 'C4':C4_deseri}
    ct_deseri = {'rkc':maabect_deseri,'ec':h['cipher']['ec']}

    hash_deseri = {'h': h_deseri, 'r': r_deseri, 'cipher':ct_deseri, 'N1': N1_deseri, 'e': e_deseri}
    return hash_deseri

def hash(msg):
    xi = chamHash.hash(pk, sk, msg)
    etd = [xi['p1'],xi['q1']]
    #if debug: print("Hash...")
    #if debug: print("hash result =>", xi)
 
    # encrypt
    rand_key = groupObj.random(GT)
    #if debug: print("msg =>", rand_key)
    #encrypt rand_key
    maabect = maabe.encrypt(public_parameters, maabepk, rand_key, access_policy)
    #rand_key->symkey AE  
    symcrypt = AuthenticatedCryptoAbstraction(extractor(rand_key))
    #symcrypt msg(etd=(p1,q1))
    etdtostr = [str(i) for i in etd]
    etdsumstr = etdtostr[0]+etdtostr[1]
    symct = symcrypt.encrypt(etdsumstr)
    ct = {'rkc':maabect,'ec':symct}
    h = {'h': xi['h'], 'r': xi['r'], 'cipher':ct, 'N1': xi['N1'], 'e': xi['e']}
    hash_seri = serializehash(h)

    return hash_seri

def check(msg, h):
    hash_deseri = deserializehash(h)
    checkresult = chamHash.hashcheck(pk, msg, hash_deseri)
    return checkresult

def collision(msg1, msg2, h):
    #deserialize h
    hash_deseri = deserializehash(h)
    #decrypt rand_key
    rec_key = maabe.decrypt(public_parameters, user_sk, hash_deseri['cipher']['rkc'])
    #rec_key->symkey AE
    rec_symcrypt = AuthenticatedCryptoAbstraction(extractor(rec_key))
    #symdecrypt rec_etdsumstr
    rec_etdsumbytes = rec_symcrypt.decrypt(hash_deseri['cipher']['ec'])
    rec_etdsumstr = str(rec_etdsumbytes, encoding="utf8")
    #print("etdsumstr type=>",type(rec_etdsumstr))
    #sumstr->etd str list
    rec_etdtolist = cut_text(rec_etdsumstr, 309)
   # print("rec_etdtolist=>",rec_etdtolist)
    #etd str list->etd integer list
    rec_etdint = {'p1': integer(int(rec_etdtolist[0])),'q1':integer(int(rec_etdtolist[1]))}
    #print("rec_etdint=>",rec_etdint)
    r1 = chamHash.collision(msg1, msg2, hash_deseri, rec_etdint, pk)
    #if debug: print("new randomness =>", r1)
    new_h = {'h': hash_deseri['h'], 'r': r1, 'cipher': hash_deseri['cipher'], 'N1': hash_deseri['N1'], 'e': hash_deseri['e']}
    newh_seri = serializehash(new_h)
    return newh_seri


def main():
    # hash
    msg = "Video provides a powerful way to help you prove your point. When you click Online Video, you can paste in the embed code for t"
    h = hash(msg)
    print("h =>", h)

    hbytes = json.dumps(h)
    print("cham Hash length = ",len(hbytes))

    # hashcheck
    checkresult = check(msg, h)
    print("checkresult =>", checkresult)

    #collision
    msg1 = "Video provides a powerful way to help you prove your point. When you click Online Video, you can paste in the embed code for p"
    new_h = collision(msg,msg1,h)
    print("new_h =>", new_h)

    checkresult2 = check(msg1, new_h)
    print("checkresult2 =>", checkresult2)
    if checkresult2: 
        print("collision generated correctly!!!")

if __name__ == '__main__':
    debug = True
    main()