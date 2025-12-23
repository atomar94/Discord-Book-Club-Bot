import discord
from discord import ui

import googlemaps
from geopy.distance import geodesic
import requests

import json
import os
import re
import logging

# Setup Google Maps client
#gmaps = googlemaps.Client(key=os.getenv("GOOGLE_MAPS_API_KEY"))
LB_COORDS = (33.7701, -118.1937) # Center of Long Beach

class CafeNominationModal(ui.Modal, title='Nominate a New Cafe'):
    # The individual text inputs in the form
    cafe_name = ui.TextInput(
        label='Cafe Name',
        placeholder='e.g., The Roasted Bean',
        min_length=2,
        max_length=50
    )
    
    maps_link = ui.TextInput(
        label='Google Maps Link',
        placeholder='https://goo.gl/maps/...',
        style=discord.TextStyle.short
    )
    
    description = ui.TextInput(
        label='Description',
        placeholder='Tell us why it is good for the book club...',
        style=discord.TextStyle.long,
        max_length=300,
        required=False
    )

    async def on_submit(self, interaction: discord.Interaction):
        await interaction.response.defer(ephemeral=True)
        user_url = self.maps_link.value
        cafe_name = self.cafe_name.value

        try:
            # 1. Unshorten the URL
            session = requests.Session()
            resp = session.head(user_url, allow_redirects=True, timeout=5)
            long_url = resp.url

            # 2. Extract Latitude/Longitude from the URL
            # Look for the pattern @lat,lng (e.g., @33.77,-118.19)
            coord_match = re.search(r'@(-?\d+\.\d+),(-?\d+\.\d+)', long_url)
            
            if not coord_match:
                return await interaction.followup.send("I couldn't parse this Google Maps URL", ephemeral=True)

            lat = float(coord_match.group(1))
            lng = float(coord_match.group(2))
            found_coords = (lat, lng)

            # 3. Distance Check
            dist = geodesic(LB_COORDS, found_coords).miles
            if dist > 10:
                return await interaction.followup.send(f"❌ This cafe is a bit too far!", ephemeral=True)

            # Success! Create the voting embed
            preamble = "Want to add this cafe to the schedule? React with a ❤️"

            embed = discord.Embed(title=f"☕ Cafe Nomination: {cafe_name}", description=preamble, color=0x3498db)

            embed.add_field(name="Cafe", value=cafe_name)

            embed.add_field(name="Directions", value=f"[Google Maps]({user_url})")

            about_slug = ""
            if self.description.value != "":
                about_slug += self.description.value

            embed.add_field(name="☕ About this cafe", value=about_slug)
            embed.url = user_url
            
            vote_msg = await interaction.channel.send(embed=embed)
            await interaction.followup.send(
                f"Thanks for nominating a cafe! If we get enough (3) ❤️ reacts we'll add it to our schedule. We think {cafe_name} looks great!",
                ephemeral=True)
            await vote_msg.add_reaction("❤️")

        except Exception as e:
            logging.error(e.with_traceback())
            await interaction.followup.send(f"Error parsing link: {e}", ephemeral=True)

            
            # Post to the public channel for voting
            channel = interaction.channel
            vote_msg = await channel.send(embed=embed)
            await vote_msg.add_reaction("❤️")