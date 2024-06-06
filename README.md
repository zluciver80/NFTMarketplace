# SolanaNFTMarketplace

This project is a decentralized NFT marketplace that allows users to create, buy, and sell NFTs using Solana for the blockchain backend, Go for the server, and HTML/CSS/JavaScript for the frontend.

## Structure

- **HTML/CSS/JavaScript**: Handles the frontend interface.
- **Go**: Manages backend processing and database interactions.
- **Solana**: Manages decentralized NFT transactions on the Solana blockchain.

## Setup

### HTML/CSS/JavaScript
1. Navigate to the `root` directory.
2. Install the required dependencies:
    ```
    npm install
    ```
3. Start the development server:
    ```
    npm start
    ```

### Go
1. Install Go if it is not already installed.
2. Navigate to the project directory and set up the environment variables by creating a `.env` file in the `root` directory with the following content:
    ```
    DATABASE_URL=postgres
