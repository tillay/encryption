if ! command -v wl-copy &>/dev/null && ! command -v xclip &>/dev/null; then
  echo "No clipboard utility (wl-copy or xclip) detected. Please install one for your desktop environment."
  echo "Only the backup option is available."
  python3 minimal.py
  exit 1
fi
command -v go &>/dev/null || disable_pgp=true
python3 -c "import Crypto" &>/dev/null || pycryptodome_installed=false
echo "Choose an option:"
[[ $pycryptodome_installed != false ]] && {
  echo "1. Encrypt/Decrypt text"
  echo "2. Encrypt/Decrypt image"
  echo "3. New random password"
  echo "4. New custom password"
  echo "5. Backup script"
} || echo "please install pycryptodome library for more options"
[[ -z $disable_pgp ]] && echo "6. Open PGP menu"
read -r choice
case $choice in
  1) python3 encrypt.py ;;
  2) python3 image.py ;;
  3) 
    echo -n "Password length: "
    read -r length
    python3 encrypt.py -p "$length"
    ;;
  4) python3 encrypt.py -n ;;
  5) python3 minimal.py ;;
  6) 
    [[ ! -f ./pgpcli/bin/pgpcli ]] && {
      mkdir -p ./pgpcli/bin
      (cd pgpcli && go build -o bin/pgpcli)
    }
    ./pgpcli/bin/pgpcli
    ;;
  *) echo "Invalid option, skill issue detected." ;;
esac
