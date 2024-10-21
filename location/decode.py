import base64
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives.padding import PKCS7
import codecs
import struct
import hashlib

payload = "LMC3CQMEuRm81M3E0+Y70gb8KL1NEZ1rdj1GWz7RzlZR0N3odV6BqhJFW0CizXrrKCZk3+/g26hvgDJF6YAfdqfbYsMIHj4VoCuSFFMAWLdbavJuTK2t7w=="

data = base64.b64decode(payload)


def bytes_to_int(b):
    return int(codecs.encode(b, "hex"), 16)


def sha256(data):
    digest = hashlib.new("sha256")
    digest.update(data)
    return digest.digest()


def decrypt(enc_data, algorithm_dkey, mode):
    decryptor = Cipher(algorithm_dkey, mode, default_backend()).decryptor()
    return decryptor.update(enc_data) + decryptor.finalize()


def unpad(paddedBinary, blocksize):
    unpadder = PKCS7(blocksize).unpadder()
    return unpadder.update(paddedBinary) + unpadder.finalize()


def decode_tag(data):
    latitude = struct.unpack(">i", data[0:4])[0] / 10000000.0
    longitude = struct.unpack(">i", data[4:8])[0] / 10000000.0
    confidence = bytes_to_int(data[8:9])
    status = bytes_to_int(data[9:10])
    return {"latitude": latitude, "longitude": longitude, "confidence": confidence, "status": status}


timestamp = bytes_to_int(data[0:4])
priv = bytes_to_int(base64.b64decode(
    "63wf6z/O7aasxWSD64I48IK/wROwBSDxeUjiJw=="))
eph_key = ec.EllipticCurvePublicKey.from_encoded_point(
    ec.SECP224R1(), data[5:62])
shared_key = ec.derive_private_key(priv, ec.SECP224R1(), default_backend()).exchange(
    ec.ECDH(), eph_key
)
symmetric_key = sha256(shared_key + b"\x00\x00\x00\x01" + data[5:62])
decryption_key = symmetric_key[:16]
iv = symmetric_key[16:]
enc_data = data[62:72]
tag = data[72:]

decrypted = decrypt(enc_data, algorithms.AES(
    decryption_key), modes.GCM(iv, tag))
res = decode_tag(decrypted)

print(res)
