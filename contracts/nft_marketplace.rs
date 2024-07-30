use anchor_lang::prelude::*;
use anchor_lang::solana_program::{program::invoke, pubkey::Pubkey, system_instruction};

declare_id!("Fg6PaFyLnbtBRQdCDVLyNtQp4A6XYj9oHKtDk1y1s6Zn");

#[program]
pub mod solana_nft_marketplace {
    use super::*;

    pub fn create_nft(
        ctx: Context<CreateNft>,
        title: String,
        uri: String,
    ) -> Result<()> {
        let nft_account = &mut ctx.accounts.nft_account;
        nft_account.owner = *ctx.accounts.user.key;
        nft_account.title = title;
        nft_account.uri = uri;
        nft_account.listed = false;
        nft_account.price = 0;
        Ok(())
    }

    pub fn list_nft_for_sale(ctx: Context<ListNftForSale>, price: u64) -> Result<()> {
        let nft_account = &mut ctx.accounts.nft_account;
        require!(price > 0, ErrorCode::PriceMustBeAboveZero);
        require!(!nft_account.listed, ErrorCode::NftAlreadyListed);
        nft_account.listed = true;
        nft_account.price = price;
        Ok(())
    }

    pub fn buy_nft(ctx: Context<BuyNft>, expected_price: u64) -> Result<()> {
        let nft_account = &mut ctx.accounts.nft_account;
        
        require!(nft_account.listed, ErrorCode::NftNotForSale);
        require!(nft_account.price == expected_price, ErrorCode::IncorrectPrice);

        invoke(
            &system_instruction::transfer(
                &ctx.accounts.buyer.key,
                &nft_account.owner,
                nft_account.price,
            ),
            &[
                ctx.accounts.buyer.to_account_info(),
                ctx.accounts.nft_owner.to_account_info(),
                ctx.accounts.system_program.to_account_info(),
            ],
        )?;

        nft_account.owner = *ctx.accounts.buyer.key;
        nft_account.listed = false;
        nft_account.price = 0;

        Ok(())
    }

    pub fn get_nft_details(_ctx: Context<GetNftDetails>) -> Result<()> {
        Ok(())
    }
}

#[derive(Accounts)]
pub struct CreateNft<'info> {
    #[account(init, payer = user, space = 1024)]
    pub nft_account: Account<'info, NftAccount>,
    #[account(mut)]
    pub user: Signer<'info>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct ListNftForSale<'info> {
    #[account(mut, has_one = owner)]
    pub nft_account: Account<'info, NftAccount>,
    pub owner: Signer<'info>,
}

#[derive(Accounts)]
pub struct BuyNft<'info> {
    #[account(mut, has_one = owner)]
    pub nft_account: Account<'info, NftAccount>,
    #[account(mut)]
    pub nft_owner: AccountInfo<'info>,
    #[account(mut)]
    pub buyer: Signer<'info>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct GetNftDetails {}

#[account]
pub struct NftAccount {
    pub owner: Pubkey,
    pub title: String,
    pub uri: String,
    pub listed: bool,
    pub price: u64,
}

#[error_code]
pub enum ErrorCode {
    #[msg("Price must be above zero")]
    PriceMustBeAboveZero,
    #[msg("NFT is already listed for sale")]
    NftAlreadyListed,
    #[msg("This NFT is not listed for sale")]
    NftNotForSale,
    #[msg("The price does not match the listing price")]
    IncorrectPrice,
}