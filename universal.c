#include <openssl/aes.h>
#include <openssl/rand.h>
#include <openssl/evp.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <openssl/bio.h>
#include <openssl/buffer.h>

#define BLOCK_SIZE AES_BLOCK_SIZE
#define KEY_SIZE 32

void deriveKey(const char *password, unsigned char *key) {
    EVP_Digest(password, strlen(password), key, NULL, EVP_sha256(), NULL);
}

char *base64Encode(const unsigned char *input, int length) {
    BIO *b64 = BIO_new(BIO_f_base64()), *bmem = BIO_new(BIO_s_mem());
    BUF_MEM *bptr;
    b64 = BIO_push(b64, bmem);
    BIO_set_flags(b64, BIO_FLAGS_BASE64_NO_NL);
    BIO_write(b64, input, length);
    BIO_flush(b64);
    BIO_get_mem_ptr(b64, &bptr);
    char *buff = malloc(bptr->length + 1);
    memcpy(buff, bptr->data, bptr->length);
    buff[bptr->length] = '\0';
    BIO_free_all(b64);
    return buff;
}

unsigned char *base64Decode(const char *input, int *length) {
    BIO *b64 = BIO_new(BIO_f_base64()), *bmem = BIO_new_mem_buf(input, -1);
    unsigned char *buffer = malloc(strlen(input));
    memset(buffer, 0, strlen(input));
    bmem = BIO_push(b64, bmem);
    BIO_set_flags(b64, BIO_FLAGS_BASE64_NO_NL);
    *length = BIO_read(b64, buffer, strlen(input));
    BIO_free_all(b64);
    return buffer;
}

void encrypt(const char *plaintext, const char *password, char *output) {
    unsigned char key[KEY_SIZE], iv[BLOCK_SIZE], ciphertext[1024];
    int len, ciphertext_len;
    deriveKey(password, key);
    RAND_bytes(iv, BLOCK_SIZE);
    EVP_CIPHER_CTX *ctx = EVP_CIPHER_CTX_new();
    EVP_EncryptInit_ex(ctx, EVP_aes_256_cbc(), NULL, key, iv);
    EVP_EncryptUpdate(ctx, ciphertext, &len, (unsigned char *)plaintext, strlen(plaintext));
    ciphertext_len = len;
    EVP_EncryptFinal_ex(ctx, ciphertext + len, &len);
    ciphertext_len += len;
    EVP_CIPHER_CTX_free(ctx);
    unsigned char combined[BLOCK_SIZE + ciphertext_len];
    memcpy(combined, iv, BLOCK_SIZE);
    memcpy(combined + BLOCK_SIZE, ciphertext, ciphertext_len);
    char *encoded = base64Encode(combined, BLOCK_SIZE + ciphertext_len);
    strcpy(output, encoded);
    free(encoded);
}

void decrypt(const char *ciphertext, const char *password, char *output) {
    unsigned char key[KEY_SIZE], iv[BLOCK_SIZE], plaintext[1024], *decoded;
    int len, plaintext_len, decoded_len;
    deriveKey(password, key);
    decoded = base64Decode(ciphertext, &decoded_len);
    memcpy(iv, decoded, BLOCK_SIZE);
    EVP_CIPHER_CTX *ctx = EVP_CIPHER_CTX_new();
    EVP_DecryptInit_ex(ctx, EVP_aes_256_cbc(), NULL, key, iv);
    EVP_DecryptUpdate(ctx, plaintext, &len, decoded + BLOCK_SIZE, decoded_len - BLOCK_SIZE);
    plaintext_len = len;
    EVP_DecryptFinal_ex(ctx, plaintext + len, &len);
    plaintext_len += len;
    EVP_CIPHER_CTX_free(ctx);
    free(decoded);
    plaintext[plaintext_len] = '\0';
    strcpy(output, (char *)plaintext);
}
int main(int argc, char *argv[]) {
    if (argc != 2) {
        printf("Usage: %s <password>\n", argv[0]);
        return 1;
    }

    char *password = argv[1];
    char input[1024], encrypted[2048], decrypted[1024];
    printf("Input text: ");
    fgets(input, sizeof(input), stdin);
    input[strcspn(input, "\n")] = 0;

    if (strncmp(input, "&&", 2) == 0) {
        decrypt(input + 2, password, decrypted);
        printf("Decrypted text: %s\n", decrypted);
    } else {
        encrypt(input, password, encrypted);
        printf("Encrypted text: &&%s\n", encrypted);
    }

    return 0;
}

