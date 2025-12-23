import discord
from discord import ui

from schedule.schedule import push_book

class BookRecommendationModal(ui.Modal, title='Nominate a Book'):
    # Title is required
    book_title = ui.TextInput(
        label='Book Title',
        placeholder='e.g., The Great Gatsby',
        min_length=1,
        max_length=100,
        required=True
    )
    
    # Goodreads link is optional
    goodreads_link = ui.TextInput(
        label='Goodreads Link (Optional)',
        placeholder='https://www.goodreads.com/book/show/...',
        required=False
    )
    
    # Description is optional and uses a larger text box
    description = ui.TextInput(
        label='Why should we read this? (Optional)',
        style=discord.TextStyle.long,
        placeholder='Briefly describe the plot or why you liked it...',
        required=False,
        max_length=400
    )

    async def on_submit(self, interaction: discord.Interaction):
        # Create a clean Embed to display the recommendation
        embed = discord.Embed(
            title=f"📖 Book Recommendation: {self.book_title.value}",
            color=discord.Color.green(),
            timestamp=interaction.created_at
        )
        
        embed.add_field(name="Title", value=self.book_title.value)

        # Add the Goodreads link if provided
        if self.goodreads_link:
            link_val = self.goodreads_link.value if self.goodreads_link.value else "No link provided."
            embed.add_field(name="Link", value=link_val, inline=False)
        
        # Add the description if provided
        if self.description:
            desc_val = self.description.value if self.description.value else "No description provided."
            embed.add_field(name="Description", value=desc_val, inline=False)
        
        embed.set_footer(text=f"Nominated by {interaction.user.display_name}")

        push_book(self.book_title.value, "", "")

        # Post the nomination to the channel
        await interaction.response.send_message(f"Nomination for *{self.book_title.value}* submitted!", ephemeral=True)
        
        vote_msg = await interaction.channel.send(embed=embed)
        # Add a heart for voting, just like the cafes
        await vote_msg.add_reaction("❤️")