# Use the latest Python image as the base
FROM python:latest

# Set the working directory
WORKDIR /usr/src/app

# Copy the requirements file and install Python dependencies
COPY requirements.txt ./
RUN pip install --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt

# Update and install necessary packages
RUN apt-get update && \
    apt-get install -y wget gnupg unzip && \
    apt-get install -y firefox-esr \
                       ca-certificates \
                       fonts-liberation \
                       libappindicator3-1 \
                       libasound2 \
                       libatk-bridge2.0-0 \
                       libatk1.0-0 \
                       libcups2 \
                       libdbus-1-3 \
                       libdrm2 \
                       libxcomposite1 \
                       libxdamage1 \
                       libxrandr2 \
                       libgbm1 \
                       libgtk-3-0 \
                       libnss3 \
                       libxss1 \
                       libxtst6 \
                       lsb-release \
                       xdg-utils \
                       xauth \
                       xvfb

# Install Google Chrome
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
    sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list' && \
    apt-get update && \
    apt-get install -y google-chrome-stable

# Install ChromeDriver
RUN wget -O /tmp/chromedriver.zip https://storage.googleapis.com/chrome-for-testing-public/126.0.6478.55/linux64/chromedriver-linux64.zip && \
    unzip /tmp/chromedriver.zip -d /tmp/ && \
    mv /tmp/chromedriver-linux64/chromedriver /usr/local/bin/chromedriver && \
    rm -rf /tmp/chromedriver.zip /tmp/chromedriver-linux64

# Set the command to run Robot Framework tests
CMD ["robot", "tests/tests.robot"]
