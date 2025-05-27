from flask import Flask, request, jsonify
import crypto, os, base64, json

# Пока тестируем без ECDH, ключи будут заданы заранее, но по умолчанию None
server_aes_key = b'\xb9M\x0b8\x00\x10\x90\x16\xc7\xed\x93\x08\xc1\x00J\xf2\x08\xb0\x01~\xb5_G\x805\xac\x95\xa2t`1\xde'
server_hmac_key = b'Dp\xc2\xc6B\x16\xb8\\\xaf_z5\x8dC\x1f3\x19\n\xe1u8\xe1Q:\xd1}\xb2\xa0\xf8$\xa6\x0e'

PORT = 7654
HOST = '0.0.0.0'

app = Flask(__name__)

# Json такого вида:
# {
#     "content": [
#         "base64string"
#         ]
# }

@app.route('/encrypt', methods=['POST'])
def encrypt():
    data = request.get_json()
    content = [base64.b64decode(i).decode("utf-8") for i in data['content']]

    json_str = json.dumps(content)
    content = json_str.encode('utf-8')

    print(f"\nserver.py | encrypt() json_str: {json_str}\n")
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

    return jsonify(encrypted_text)


@app.route('/decrypt', methods=['POST'])
def decrypt():
    data = request.get_json()
    # content = base64.b64decode(data['content'])
    print(f"\nserver.py | decrypt() data: {data}\n")
    # print(content)
    # print(base64.b64decode(data["iv"]))
    # print("\n\n\n\n\n\n")

    crypto_box = crypto.Aes256CbcHmac(server_aes_key, server_hmac_key)

    decrypted_text = crypto_box.decrypt(data)

    data_list = json.loads(decrypted_text.decode('utf-8'))
    print(f"server.py | decrypt() result: 'content': {data_list}")
    content_base64_list = [base64.b64encode(item.encode('utf-8')).decode('utf-8')
        for key, item in data_list.items()]
    result = {
        "content": content_base64_list
    }

    # return jsonify(base64.b64encode(decrypted_text).decode())
    print(f"\nserver.py | decrypt() result: {result}\n")
    return result

@app.route('/ecdh', methods=['POST'])
def ecdh():
    # {
    #     content: ["base64string_client_public_key"]
    # }
    data = request.get_json()
    print(f"\nserver.py | ecdh() data: {data}\n")
    # json_str = json.dumps(data["content"])
    # print(f"\nserver.py | ecdh() json_str: {json_str}\n")
    client_pub = json.dumps(data["content"])
    print(f"\nserver.py | ecdh() client_pub: {client_pub}\n")

    server = crypto.ECDHKeyExchange() # 4
    server_pub = server.get_public_key_base64() # 5

    server.compute_shared_secret(client_pub) # 7,9

    # Ключи снизу используем для шифрования и проверки целостности
    aes_key = server.aes_key
    hmac_key = server.hmac_key

    result = {
        "content": server_pub
    }

    return result

if __name__ == '__main__':
    app.run(host=HOST, port=PORT, debug=True)



# JS скрипты для запросов через Chrome DevTools
# btoa(str) - перевод str в b64 в JS

# fetch('http://127.0.0.1:7654/encrypt', {
#     method: 'POST',
#     headers: {
#         'Content-Type': 'application/json'
#     },
#     body: JSON.stringify({
#         content: [btoa("testing bebra"), btoa("asdfasdf"), btoa("2"), btoa("3"), btoa("4"), "0LPQvtC50LTQsA=="]
#     })
# })
# .then(response => response.json())
# .then(data => console.log(data))
# .catch(error => console.error(error));


# fetch('http://127.0.0.1:7654/decrypt', {
#     method: 'POST',
#     headers: {
#         'Content-Type': 'application/json'
#     },
#     body: JSON.stringify({content: '6J87GLgA4jpP3DtPQCX2/+rXctJJUmECYj58Jari+xcDUUMOmrcd+ZnZZ7HmtL7d', hmac: 'YVlgDYFz1CVlDaDGvhP2gqB5IhwfyGbV4iRo5+gwrPQ=', iv: 'oFtTJ+F57vo9nqqhmo3y2Q==', nonce: 'Nf/DN0FGfUM0bDoz'})
# })
# .then(response => response.json())
# .then(data => console.log(data))
# .catch(error => console.error(error));