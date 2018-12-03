# Burst Wallet

## Key Features
- 🔒 PIN-based access, no passphrases stored
- 🔥 Latest Angular, TypeScript, LokiJS
- 🔥 [angular-cli](https://cli.angular.io/) support
- 🔥 [Hot Module Replacement](https://webpack.js.org/concepts/hot-module-replacement/) support

## Contribute / Development Setup

1. Install Node JS and npm
    * Windows
        * https://nodejs.org/en/download/
    * Ubuntu/Debian
        * `curl -sL https://deb.nodesource.com/setup_7.x | sudo -E bash -`
        * `sudo apt-get install -y nodejs`
        * `sudo apt-get install -y build-essential`

2. Clone this project:
    ```
    git clone git@github.com:PoC-Consortium/Aspera.git
    cd Aspera/web
    ```

3. Uninstall old versions of the cli and Install latest version of angular-cli
    ```
    npm uninstall -g angular-cli @angular/cli
    npm cache clean
    npm install -g @angular/cli@latest
    ```

3. Install node modules for the project
    ```
    npm install
    ```

4. Start Webpack Server

    ```
    npm run start
    ```

5. Navigate with your browser to localhost:4200 to see the UI
