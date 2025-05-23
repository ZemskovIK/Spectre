import base64
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives.kdf.hkdf import HKDF
from cryptography.hazmat.primitives import hashes, serialization


class ECDHKeyExchange:
    def __init__(self):
        self._private_key = ec.generate_private_key(ec.SECP256R1())
        self._shared_secret = None
        self.aes_key = None
        self.hmac_key = None

    def get_public_key_base64(self) -> str:
        # Возвращает публичный ключ в base64
        public_bytes = self._private_key.public_key().public_bytes(
            encoding=serialization.Encoding.X962,
            format=serialization.PublicFormat.UncompressedPoint
        )
        return base64.b64encode(public_bytes).decode()

    def compute_shared_secret(self, other_pub_base64: str):
        # Вычисляет ключи AES + HMAC по чужому публичному ключу
        other_pub_bytes = base64.b64decode(other_pub_base64)
        other_public_key = ec.EllipticCurvePublicKey.from_encoded_point(
            ec.SECP256R1(), other_pub_bytes
        )
        self._shared_secret = self._private_key.exchange(ec.ECDH(), other_public_key)

        derived = HKDF(
            algorithm=hashes.SHA256(),
            length=64,
            salt=None,
            info=b"handshake data"
        ).derive(self._shared_secret)

        self.aes_key = derived[:32]
        self.hmac_key = derived[32:]



def main():
    # Воображаемая архитектура клиент-сервер
    client = ECDHKeyExchange() # 1 
    client_pub = client.get_public_key_base64() # 2

    server = ECDHKeyExchange() # 4
    server_pub = server.get_public_key_base64() # 5


    # Воображаемый обмен ключами # 3

    server.compute_shared_secret(client_pub) # 7,9
    client.compute_shared_secret(server_pub) # 8,9

    # Проверка: ключи совпадают !!! сделать проверка на совпадение, иначе перехендшейк руками)

    if server.aes_key != client.aes_key: "AES ключи не совпадают!"
    if server.hmac_key != client.hmac_key: "HMAC ключи не совпадают!"
    print(f"{server.aes_key}\n{client.aes_key} ")


if __name__ == "__main__":
    main()