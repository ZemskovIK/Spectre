from flask import Flask, request, jsonify
import crypto, os, base64, json

# Пока тестируем без ECDH, ключи будут заданы заранее, но по умолчанию None
server_aes_key = b'\xb9M\x0b8\x00\x10\x90\x16\xc7\xed\x93\x08\xc1\x00J\xf2\x08\xb0\x01~\xb5_G\x805\xac\x95\xa2t`1\xde'
server_hmac_key = b'Dp\xc2\xc6B\x16\xb8\\\xaf_z5\x8dC\x1f3\x19\n\xe1u8\xe1Q:\xd1}\xb2\xa0\xf8$\xa6\x0e'

PORT = 7654
HOST = '0.0.0.0'

app = Flask(__name__)
keys_by_users = {}

# Json такого вида:
# {
#     "content": [
#         "base64string"
#         ]
# }

@app.route('/encrypt', methods=['POST'])
def encrypt():
    data = request.get_json()
    print(f"server.py | encrypt() data: {data}, {type(data)}")
    content = [base64.b64decode(i).decode("utf-8") for i in data['content']]

    json_str = json.dumps(content)
    content = json_str.encode('utf-8')

    print(f"\nserver.py | encrypt() json_str: {json_str}, {type(json_str)}\n")
    print(f"\nserver.py | encrypt() content: {content}\n")

    crypto_box = crypto.Aes256CbcHmac(server_aes_key, server_hmac_key)
    nonce = os.urandom(12)

    encrypted_text = crypto_box.encrypt(content, nonce)

    # data_list = json.loads(encrypted_text.decode('utf-8'))
    # content_base64_list = [base64.b64encode(item.encode('utf-8')).decode('utf-8')
    #     for item in data_list]
    # result = {
    #     "content": content_base64_list
    # }
    print(f"\nserver.py | encrypt(): {(encrypted_text)}, {type((encrypted_text))}\n")
    return jsonify(encrypted_text)


@app.route('/decrypt', methods=['POST'])
def decrypt():
    data = request.get_json()
    # content = base64.b64decode(data['content'])
    print(f"\nserver.py | decrypt() data: {data}, {type(data)}\n")
    # print(content)
    # print(base64.b64decode(data["iv"]))
    # print("\n\n\n\n\n\n")
    data = json.loads(data)
    crypto_box = crypto.Aes256CbcHmac(server_aes_key, server_hmac_key)

    decrypted_text = crypto_box.decrypt(data)
    print(f"server.py | decrypt() decrypted_text: {decrypted_text}, {type(decrypted_text)}")

    data_list = json.loads(decrypted_text.decode('utf-8'))
    print(f"server.py | decrypt() 'content': {data_list}, {type(data_list)}")

    # for key, item in data_list.items():
    #     print(f"server.py | decrypt() key, item: {key}, {item}")
    
    content_base64_list = [base64.b64encode(item.encode('utf-8')).decode('utf-8')
        for item in data_list]
    
    for item in data_list:
        print(f"items: {item}")

    result = {
        "content": content_base64_list
    }

    # return jsonify(base64.b64encode(decrypted_text).decode())
    print(f"\nserver.py | decrypt() result: {result}\n")
    return result

@app.route('/ecdh', methods=['POST'])
def ecdh():
    global keys_by_user
    data = request.get_json()
    # print(f"\nserver.py | ecdh() data: {data}\n")
    # json_str = json.dumps(data["content"])
    # print(f"\nserver.py | ecdh() json_str: {json_str}\n")
    if len(data) == 1:
        user_ip = json.dumps(data["from"])

        print(f"\nserver.py | ecdh()1 user_ip: {user_ip}\n")

        server = crypto.ECDHKeyExchange()
        server_pub = server.get_public_key_base64()
        keys_by_users[user_ip] = [server_pub, server._private_key]

        print(f"\nserver.py | ecdh()1 server_pub: {server_pub}\n")

        result = {
            "key":server_pub
        }

    elif len(data) == 2:
        user_ip = json.dumps(data["from"])
        client_public_key =  json.dumps(data["key"])

        server = crypto.ECDHKeyExchange()
        server._private_key = keys_by_users[user_ip][1]

        server.compute_shared_secret(client_public_key)

        keys_by_users[user_ip] = [server.aes_key, server.hmac_key]

        return '', 204

    print(f"\nserver.py | ecdh() keys_by_users: {keys_by_users}\n")

    return result

if __name__ == '__main__':
    app.run(host=HOST, port=PORT, debug=True)