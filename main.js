import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import axios from 'axios';
import { useWallet } from '@solana/wallet-adapter-react';
import { Connection, clusterApiUrl } from '@solana/web3.js';

const API_ENDPOINT_URL = process.env.REACT_APP_API_URL;
const SOLANA_CLUSTER_NETWORK = process.env.REACT_APP_SOLANA_NETWORK;

const solanaNetworkURL = clusterApiUrl(SOLANA_CLUSTER_NETWORK);
const blockchainConnection = new Connection(solanaNetworkURL, 'confirmed');

const userWallet = useWallet();

async function fetchNFTCollection() {
    try {
        const response = await axios.get(`${API_ENDPOINT_URL}/nfts`);
        return response.data;
    } catch (error) {
        console.error('Error fetching NFT collection:', error);
    }
}

const NFTGallery = ({ nfts }) => {
    if (!nfts.length) return <p>No NFTs found.</p>;
    
    return (
        <div>
            {nfts.map((nft) => (
                <div key={nft.id}>
                    <img src={nft.image} alt={nft.name} />
                    <h3>{nft.name}</h3>
                    <p>{nft.description}</p>
                </div>
            ))}
        </div>
    );
};

const WalletConnectButton = () => {
    const { connected, connect, disconnect } = useWallet();

    return (
        <div>
            {!connected ? (
                <button onClick={connect}>Connect Wallet</button>
            ) : (
                <button onClick={disconnect}>Disconnect Wallet</button>
            )}
        </div>
    );
};

const NFTMarketplace = () => {
    const [nftCollection, setNFTCollection] = useState([]);

    useEffect(() => {
        fetchNFTCollection().then(setNFTCollection);
    }, []);

    return (
        <Router>
            <div>
                <WalletConnectButton />
                <Switch>
                    <Route path="/">
                        <NFTGallery nfts={nftCollection} />
                    </Route>
                </Switch>
            </div>
        </Router>
    );
};

export default NFTMarketplace;