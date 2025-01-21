import json, os, http.client, sys
from utils import encrypt, decrypt

PASSWORD_FILE = "/tmp/key"
TOKEN_FILE = "./token"
channel_id = 1315424206523203684

def get_pass():
    if os.path.exists(PASSWORD_FILE):
        with open(PASSWORD_FILE, 'r') as f:
            return f.read()

def get_token():
    if os.path.exists(TOKEN_FILE):
        with open(TOKEN_FILE, 'r') as f:
           return f.readline().strip()

header_data = {
    "Content-Type": "application/json",
    "User-Agent": "tilley",
    "Authorization": get_token()
}

def send_message(channel_id, message_content):
    conn = http.client.HTTPSConnection("discord.com", 443)
    message_data = json.dumps({
        "content": message_content,
        "tts": False
    })
    try:
        conn.request("POST", f"/api/v10/channels/{channel_id}/messages", message_data, header_data)
        response = conn.getresponse()

        if 199 < response.status < 300:
            print("Message sent successfully.")
        else:
            print(f"Discord aint happy")
    except TypeError as e:
        print(f"Please move your token file to {TOKEN_FILE}")
        sys.exit(1)
    finally:
        conn.close()

def listen_message(channel_id):
    conn = http.client.HTTPSConnection("discord.com", 443)
    try:
        conn.request("GET", f"/api/v10/channels/{channel_id}/messages", headers=header_data)
        response = conn.getresponse()
        if 199 < response.status < 300:
            messages = json.loads(response.read().decode())
            messages.reverse()
            prev_message = ""
            for message in messages[-1:]:
                if prev_message != message['content']:
                    if "&&" in message['content'] and message['content'] != None:
                        return(f"{message['author']['username']}: {decrypt(message['content'], get_pass())}")
                        prev_message = message['content']
        else:
            return(f"Discord aint happy: {response.status} error")
    except TypeError as e:
        print(f"Please move your token file to {TOKEN_FILE}")
        sys.exit(1)
    finally:
        conn.close()

os.system("clear")
if sys.argv[1] == "listen":
    last_var = None
    while True:
        var = listen_message(channel_id)
        if var != last_var:
            print(var)
            last_var = var

elif sys.argv[1] == "send":
    while True:
        send_message(channel_id, "&&"+encrypt(input("Message to encrypt: "),get_pass()))
