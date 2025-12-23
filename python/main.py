import discord
from discord.ext import commands, tasks

from announce.announce import make_announcement
from nominate.modal import CafeNominationModal
from recommend.modal import BookRecommendationModal
from schedule.schedule import push_schedule, read_schedule, update_book_votes, get_next_book

from datetime import datetime, timedelta
import logging
import os
import re

class CafeBot(commands.Bot):
    def __init__(self, intents = None):
        intents = intents or discord.Intents.default()
        super().__init__(command_prefix="!", intents=intents)

    async def setup_hook(self):
        # This syncs your slash commands to Discord
        await self.tree.sync()
        logging.info("Tree sync")

    @tasks.loop(hours=24)
    async def weekly_schedule_update(self):
        """Checks every day if it is Monday. If so, posts the schedule."""
        
        # 0 = Monday, 1 = Tuesday, ..., 6 = Sunday
        if datetime.now().weekday() != 6:
            pass
            # return

        ANNOUNCEMENT_CHANNEL_ID = 1451785614163709964 
        channel = self.get_channel(ANNOUNCEMENT_CHANNEL_ID)
        
        if not channel:
            return

        final_message = make_accouncement()

        # 3. Send the announcement
        await channel.send(final_message)

    async def on_raw_reaction_add(self, payload):
        logging.info(f"on_raw_reaction_add (user_id: {payload.user_id})")
        if payload.user_id == self.user.id:
            logging.debug("Matching user_id, skip")
            return
        if str(payload.emoji) != "❤️":
            logging.debug("Not a heart emoji, skip")
            return

        channel = self.get_channel(payload.channel_id)
        message = await channel.fetch_message(payload.message_id)

        if message.author == self.user and message.embeds:
            embed = message.embeds[0]
            logging.debug(f"Checking message: {embed.title}")
            if "Recommendation" in embed.title:
                reaction = discord.utils.get(message.reactions, emoji="❤️")

                book_title = ""
                for field in embed.fields:
                    if field.name == "Title":
                        book_title = field.value
                        break
                if reaction and book_title:
                    update_book_votes(book_title, reaction.count)
                    return
            if "Nomination:" in embed.title and not "(✅ Success!)" in embed.title:
                reaction = discord.utils.get(message.reactions, emoji="❤️")
                checkbox = discord.utils.get(message.reactions, emoji="✅")
                if checkbox and checkbox.count > 0 and checkbox.me:
                    logging.info("Already added. skipping")
                    return
                if reaction and reaction.count >= 2:
                    logging.info("Successful voting")
                    cafe_name = ""

                    for field in embed.fields:
                        if field.name == "Cafe":
                            cafe_name = field.value
                    
                    match = re.search(r"Nomination:\s*(.*)", embed.title)
                    if match:
                        cafe_name = match.group(1)

                    push_schedule(cafe_name, embed.url)
                    await message.add_reaction("✅")
                    await channel.send(f"✅ **{cafe_name}** looks great! It has been added to the schedule. We can't wait to check it out!")

intents = discord.Intents.default()
intents.message_content = True
intents.reactions = True
bot = CafeBot(intents=intents)


@bot.tree.command(name="nominate-a-cafe", description="Suggest a new cafe for the book club")
async def nominate(interaction: discord.Interaction):
    logging.info("nominate-a-cafe")
    # Sending the modal to the user
    await interaction.response.send_modal(CafeNominationModal())


@bot.tree.command(name="recommend-a-book", description="Suggest a new cafe for the book club")
async def recommend(interaction: discord.Interaction):
    logging.info("recommend-a-book")
    # Sending the modal to the user
    await interaction.response.send_modal(BookRecommendationModal())


@bot.tree.command(name="schedule", description="See our upcoming meeting locations")
async def schedule(interaction: discord.Interaction):
    """Displays the formatted book club schedule."""
    logging.info("schedule")
    # 1. Get the structured data
    full_schedule = read_schedule()
    
    if not full_schedule:
        await interaction.response.send_message("📭 The cafe queue is currently empty. Use `/nominate` to add some!")
        return

    # 2. Build the formatted string
    # We'll use a Code Block or an Embed for better readability
    response_lines = ["🗓️ **Upcoming Book Club Schedule**", "---"]
    
    for i, item in enumerate(full_schedule):
        # We'll highlight the first one as "Next Up"
        prefix = "☕ **Next Up:**" if i == 0 else f"{i+1}."
        line = f"{prefix} `{item['date']}` — **{item['cafe']}** - [Google Maps]({item['url']})"
        response_lines.append(line)

    # 3. Join and send
    # Discord has a 2000 character limit, but a queue of ~20 cafes will fit easily
    final_message = "\n".join(response_lines)
    await interaction.response.send_message(final_message)


@bot.tree.command(name="next-book", description="DEBUG: see next book")
async def next_book(interaction: discord.Interaction):
    """Run the next book algorithm through discord."""
    logging.info("next-book")

    book = get_next_book()

    await interaction.response.send_message(f"Next book is {book['name']}")

@bot.tree.command(name="make-annoucement", description="DEBUG: see annoucement text")
async def next_book(interaction: discord.Interaction):
    """Run the next book algorithm through discord."""
    logging.info("make-annoucement")
    await interaction.response.send_message(make_announcement(dryrun=False))


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG, format='%(asctime)s - %(levelname)s - %(message)s')

    bot.run(os.getenv("DISCORD_BOT_TOKEN"))