import { Connection, PublicKey, clusterApiUrl } from '@solana/web3.js';
import { Program, AnchorProvider, web3, utils } from '@project-serum/anchor';
import { PhantomWalletAdapter } from '@solana/wallet-adapter-wallets';
import { useWallet, WalletProvider } from '@solana/wallet-adapter-react';
import { WalletAdapterNetwork } from '@solana/wallet-adapter-base';
import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Switch, Link } from 'react-router-dom';
import axios from 'axios';

const API_ENDPOINT_URL = process.env.REACT_APP_API_URL;
const SOLANA_CLUSTER_NETWORK = process.env.REACT_APP_SOLANA_DETAIL_NETWORK;

const solanaNetworkURL = clusterApiUrl(SOLANA_CLUSTER_NETWORK);
const solanaConnection = new Connection(solanaNetworkURL, 'confirmed');

const userWallet = useWallet();

async function getNFTsFromAPI() {
    try {
        const response = await axios.get(`${API_ENDPOINTUR}/nfts`);
        return response.data;
    } catch (error) {
        console.error('Failed to fetch NFT data:', error);
    }
}

const NFTDisplay = ({ nfts }) => {
    if (!nfts.length) return <p>No NFTs found</p>;
    
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

const WalletAuthentication = () => {
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

const NFTMarketplaceApp = () => {
    const [nftCollection, setNftCollection] = useState([]);

    useEffect(() => {
        getNFTsFromAPI().then(setNftCollection);
    }, []);

    return (
        <Router>
            <div>
                <WalletAuthentication />
                <Switch>
                    <Route path="/">
                        <NFTDisplay nfts={nftCollection} />
                    </Route>
                </Switch>
            </div>
        </Router>
    );
};

export default NFTMarketplaceApp;