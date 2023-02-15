export CHROME_VERSION=110.0.5481.100-1

wget --no-verbose -O /tmp/chrome.deb https://dl.google.com/linux/chrome/deb/pool/main/g/google-chrome-stable/google-chrome-stable_${CHROME_VERSION}_amd64.deb \
  && apt install -y /tmp/chrome.deb \
  && rm /tmp/chrome.deb

# Chrome driver
wget https://chromedriver.storage.googleapis.com/111.0.5563.19/chromedriver_linux64.zip
unzip chromedriver_linux64.zip
rm chromedriver_lin*