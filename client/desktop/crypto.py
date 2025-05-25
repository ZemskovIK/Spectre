import base64
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives.kdf.hkdf import HKDF
import os
import hmac as hmac_std
from cryptography.hazmat.primitives import hashes, serialization, padding, hmac
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes



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



class Aes256CbcHmac:
    def __init__(self, aes_key: bytes, hmac_key: bytes):
        self.aes_key = aes_key
        self.hmac_key = hmac_key

    def encrypt(self, plaintext: bytes, nonce: bytes) -> dict[str, str]:
        iv = os.urandom(16)

        padder = padding.PKCS7(128).padder()
        padded_data = padder.update(plaintext) + padder.finalize()

        cipher = Cipher(algorithms.AES(self.aes_key), modes.CBC(iv))
        encryptor = cipher.encryptor()
        ciphertext = encryptor.update(padded_data) + encryptor.finalize()

        h = hmac.HMAC(self.hmac_key, hashes.SHA256())
        h.update(iv + ciphertext + nonce)
        tag = h.finalize()

        return {
            "iv": base64.b64encode(iv).decode(),
            "content": base64.b64encode(ciphertext).decode(),
            "nonce": base64.b64encode(nonce).decode(),
            "hmac": base64.b64encode(tag).decode()
        }

    def decrypt(self, data: dict[str, str]) -> bytes:
        iv = base64.b64decode(data["iv"])
        ciphertext = base64.b64decode(data["content"])
        
        nonce = base64.b64decode(data["nonce"])
        tag = base64.b64decode(data["hmac"])

        print(f"\ncrypto.py | decrypt() data: {data}\n")

        h = hmac.HMAC(self.hmac_key, hashes.SHA256())
        h.update(iv + ciphertext + nonce)

        try:
            h.verify(tag)
        except Exception:
            raise ValueError("HMAC verification failed — data may be tampered")

        cipher = Cipher(algorithms.AES(self.aes_key), modes.CBC(iv))
        decryptor = cipher.decryptor()
        padded = decryptor.update(ciphertext) + decryptor.finalize()

        unpadder = padding.PKCS7(128).unpadder()
        plaintext = unpadder.update(padded) + unpadder.finalize()

        return plaintext

def main():
    print("Этот файл не должен запускаться, это модуль!")

if __name__ == "__main__":
    main()