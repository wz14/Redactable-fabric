import MAABE
import pickle
from charm.toolbox.pairinggroup import PairingGroup, GT

def merge_dicts(*dict_args):
    """
    Given any number of dicts, shallow copy and merge into a new dict,
    precedence goes to key value pairs in latter dicts.
    """
    result = {}
    for dictionary in dict_args:
        result.update(dictionary)
    return result


def main():
    groupObj = PairingGroup('SS512')

    # init object
    maabe = MAABE.MaabeRW15(groupObj)

    # generate pp
    public_parameters = maabe.setup()
    #print("public_parameters=>",public_parameters)
    g1_seri = groupObj.serialize(public_parameters['g1'])
    g2_seri = groupObj.serialize(public_parameters['g2'])
    egg_seri = groupObj.serialize(public_parameters['egg'])
    pp_seri = {'g1': g1_seri, 'g2':g2_seri, 'egg':egg_seri}

    pp_output = open("pp_output.txt", "wb")
    pickle.dump(pp_seri, pp_output)

    #generate auth key
    (pk1, sk1) = maabe.authsetup(public_parameters, 'UT')
    egga1_seri = groupObj.serialize(pk1['egga'])
    gy1_seri = groupObj.serialize(pk1['gy'])
    alpha1_seri = groupObj.serialize(sk1['alpha'])
    y1_seri = groupObj.serialize(sk1['y'])
    pk1_seri = {'name':pk1['name'], 'egga':egga1_seri, 'gy':gy1_seri}
    sk1_seri = {'name':sk1['name'], 'alpha':alpha1_seri, 'y':y1_seri}

    (pk2, sk2) = maabe.authsetup(public_parameters, 'OU')
    egga2_seri = groupObj.serialize(pk2['egga'])
    gy2_seri = groupObj.serialize(pk2['gy'])
    alpha2_seri = groupObj.serialize(sk2['alpha'])
    y2_seri = groupObj.serialize(sk2['y'])
    pk2_seri = {'name':pk2['name'], 'egga':egga2_seri, 'gy':gy2_seri}
    sk2_seri = {'name':sk2['name'], 'alpha':alpha2_seri, 'y':y2_seri}

    maabepk = {'UT': pk1_seri, 'OU': pk2_seri}
    maabesk = {'UT': sk1_seri, 'OU': sk2_seri}
    maabepk_output = open("maabepk_output.txt", "wb")
    maabesk_output = open("maabesk_output.txt", "wb")
    pickle.dump(maabepk, maabepk_output)
    pickle.dump(maabesk, maabesk_output)
    print("maabepk=>", maabepk)
    print("maabesk=>", maabesk)
    
    # generate bob attribute keys
    gid = "bob"
    user_attr1 = ['STUDENT@UT']
    user_attr2 = ['STUDENT@OU']
    user_sk1 = maabe.multiple_attributes_keygen(public_parameters, sk1, gid, user_attr1)
    user_sk2 = maabe.multiple_attributes_keygen(public_parameters, sk2, gid, user_attr2)
    K1_seri = groupObj.serialize(user_sk1['STUDENT@UT']['K'])
    KP1_seri = groupObj.serialize(user_sk1['STUDENT@UT']['KP'])
    K2_seri = groupObj.serialize(user_sk2['STUDENT@OU']['K'])
    KP2_seri = groupObj.serialize(user_sk2['STUDENT@OU']['KP'])
    usersk1_seri = {'STUDENT@UT':{'K':K1_seri,'KP':KP1_seri}}
    usersk2_seri = {'STUDENT@OU':{'K':K2_seri,'KP':KP2_seri}}

    print("user_sk1=>", user_sk1)
    print("user_sk2=>", user_sk2)
    user_sk = {'GID': gid, 'keys': merge_dicts(usersk1_seri, usersk2_seri)}
    usersk_output = open("usersk_output.txt", "wb")
    pickle.dump(user_sk, usersk_output)

if __name__ == '__main__':
    main()