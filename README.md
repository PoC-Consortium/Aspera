# Burst Particle

# Contribute / Development Setup

1. Install Node JS and npm
    * Windows
        * https://nodejs.org/en/download/
    * Ubuntu/Debian
        * `curl -sL https://deb.nodesource.com/setup_7.x | sudo -E bash -`
        * `sudo apt-get install -y nodejs`
        * `sudo apt-get install -y build-essential`

2. Clone this project:
    ```
    git clone git@github.com:cgebe/burst-particle.git
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
    ng serve
    ```

5. Navigate with your browser to localhost:4200 to see the UI
