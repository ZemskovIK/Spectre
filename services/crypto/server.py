from flask import Flask, request, jsonify
import crypto, os, base64

# Пока тестируем без ECDH, ключи будут заданы заранее 
server_aes_key = b'\xb9M\x0b8\x00\x10\x90\x16\xc7\xed\x93\x08\xc1\x00J\xf2\x08\xb0\x01~\xb5_G\x805\xac\x95\xa2t`1\xde'
server_hmac_key = b'Dp\xc2\xc6B\x16\xb8\\\xaf_z5\x8dC\x1f3\x19\n\xe1u8\xe1Q:\xd1}\xb2\xa0\xf8$\xa6\x0e'

PORT = 7654
HOST = '0.0.0.0'

app = Flask(__name__)

# Json такого вида:
# json
# {
#     "content": [
#         "base64string"
#         ]
# }

@app.route('/encrypt_bytes', methods=['POST'])
def encrypt_bytes():
    data = request.get_json()
    content = base64.b64decode(data['content'][0])

    crypto_box = crypto.Aes256CbcHmac(server_aes_key, server_hmac_key)
    nonce = os.urandom(12)

    encrypted_text = crypto_box.encrypt(content, nonce)

    return jsonify(encrypted_text)


if __name__ == '__main__':
    app.run(host=HOST, port=PORT, debug=True)