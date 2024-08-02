use anchor_lang::prelude::*;
use anchor_lang::solana_program::{program::invoke, pubkey::Pubkey, system_instruction};

declare_id!("Fg6PaFyLnbtBRQdCDVLyNtQp4A6XYj9oHKtDk1y1s6Zn");

#[program]
pub mod solana_nft_marketplace {
    use super::*;

    pub fn create_nft_listing(
        ctx: Context<CreateNftListing>,
        nft_title: String,
        nft_uri: String,
    ) -> Result<()> {
        let listing_account = &mut ctx.accounts.listing_account;
        listing_account.owner = *ctx.accounts.creator.key;
        listing_account.title = nft_title;
        listing_account.uri = nft_uri;
        listing_account.is_listed = false;
        listing_account.sale_price = 0;
        Ok(())
    }

    pub fn list_nft(
        ctx: Context<ListForSale>,
        listing_price: u64,
    ) -> Result<()> {
        let listing_account = &mut ctx.accounts.listing_account;
        require!(listing_price > 0, ErrorCode::PriceMustBeAboveZero);
        require!(!listing_account.is_listed, ErrorCode::NftAlreadyListed);
        listing_account.is_listed = true;
        listing_account.sale_price = listing_price;
        Ok(())
    }

    pub fn purchase_nft(
        ctx: Context<PurchaseNft>,
        expected_sale_price: u64,
    ) -> Result<()> {
        let listing_account = &mut ctx.accounts.listing_account;
        
        require!(listing_account.is_listed, ErrorCode::NftNotForSale);
        require!(listing_account.sale_price == expected_sale_price, ErrorCode::IncorrectPrice);

        invoke(
            &system_instruction::transfer(
                &ctx.accounts.purchaser.key,
                &listing_account.owner,
                listing_account.sale_price,
            ),
            &[
                ctx.accounts.purchaser.to_account_info(),
                ctx.accounts.listing_owner.to_account_info(),
                ctx.accounts.system_program.to_account_info(),
            ],
        )?;

        listing_account.owner = *ctx.accounts.purchaser.key;
        listing_account.is_listed = false;
        listing_account.sale_price = 0;

        Ok(())
    }

    pub fn retrieve_nft_details(_ctx: Context<RetrieveNftDetails>) -> Result<()> {
        Ok(())
    }
}

#[derive(Accounts)]
pub struct CreateNftListing<'info> {
    #[account(init, payer = creator, space = 1024)]
    pub listing_account: Account<'info, NftListingAccount>,
    #[account(mut)]
    pub creator: Signer<'info>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct ListForSale<'info> {
    #[account(mut, has_one = owner)]
    pub listing_account: Account<'info, NftListingAccount>,
    pub owner: Signer<'info>,
}

#[derive(Accounts)]
pub struct PurchaseNft<'info> {
    #[account(mut, has_one = owner)]
    pub listing_account: Account<'info, NftListingAccount>,
    #[account(mut)]
    pub listing_owner: AccountInfo<'info>,
    #[account(mut)]
    pub purchaser: Signer<'info>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct RetrieveNftDetails {}

#[account]
pub struct NftListingAccount {
    pub owner: Pubkey,
    pub title: String,
    pub uri: String,
    pub is_listed: bool,
    pub sale_price: u64,
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