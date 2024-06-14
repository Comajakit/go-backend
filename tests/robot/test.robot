*** Settings ***
Documentation     Test Google and Bing with Selenium in headless mode.
Library           SeleniumLibrary

*** Variables ***
@{HEADLESS_OPTIONS}    headless    disable-gpu    no-sandbox    disable-dev-shm-usage

*** Keywords ***
Open Headless Chrome
    [Arguments]    ${url}
    ${options}=    Evaluate    sys.modules['selenium.webdriver'].ChromeOptions()    sys, selenium.webdriver
    FOR    ${option}    IN    @{HEADLESS_OPTIONS}
        Call Method    ${options}    add_argument    --${option}
    END
    Create Webdriver    Chrome    options=${options}
    Go To    ${url}

Open Headless Firefox
    [Arguments]    ${url}
    ${options}=    Evaluate    sys.modules['selenium.webdriver'].FirefoxOptions()    sys, selenium.webdriver
    Call Method    ${options}    add_argument    -headless
    Create Webdriver    Firefox    options=${options}
    Go To    ${url}

*** Test Cases ***
Open Google Chrome
    Open Headless Chrome    https://www.google.com
    Title Should Be    Google
    [Teardown]    Close Browser

Open Bing Chrome
    Open Headless Chrome    https://www.bing.com
    Title Should Be    Bing
    [Teardown]    Close Browser

Open Google Firefox
    Open Headless Firefox    https://www.google.com
    Title Should Be    Google
    [Teardown]    Close Browser

Open Bing Firefox
    Open Headless Firefox    https://www.bing.com
    Title Should Be    Bing
    [Teardown]    Close Browser
