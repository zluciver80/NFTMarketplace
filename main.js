import { Connection, PublicKey, clusterApiUrl } from '@solana/web3.js';
import { Program, AnchorProvider, web3, utils } from '@project-serum/anchor';
import { PhantomWalletAdapter } from '@solana/wallet-adapter-wallets';
import { useWallet, WalletProvider } from '@solana/wallet-adapter-react';
import { WalletAdapterNetwork } from '@solana/wallet-adapter-base';
import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Switch, Link } from 'react-router-dom';
import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL;
const SOLANA_NETWORK = process.env.REACT_APP_SOLANA_NETWORK;

const network = clusterApiUrl(SOLANA_NETWORK);
const connection = new Connection(network, 'confirmed');

const wallet = useWallet();

async function fetchNFTData() {
    try {
        const { data } = await axios.get(`${API_URL}/nfts`);
        return data;
    } catch (error) {
        console.error('Error fetching NFT data:', error);
    }
}

const DisplayNFTs = ({ nfts }) => {
    if (!nfts.length) return <p>No NFTs available</p>;
    
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

const Auth = () => {
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

const App = () => {
    const [nfts, setNfts] = useState([]);

    useEffect(() => {
        fetchNFTData().then(setNfts);
    }, []);

    return (
        <Router>
            <div>
                <Auth />
                <Switch>
                    <Route path="/">
                        <DisplayNFTs nfts={nfts} />
                    </Route>
                </Switch>
            </div>
        </Router>
    );
};

export default App;