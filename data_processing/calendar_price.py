import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from matplotlib.patches import Polygon

# Credit for the original calendar code: https://stackoverflow.com/a/61277350


# Settings
years = [2021, 2022]  # [2018, 2019, 2020]
weeks = [1, 2, 3, 4, 5, 6]
days = ['M', 'T', 'W', 'T', 'F', 'S', 'S']
month_names = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August',
               'September', 'October', 'November', 'December']


def split_months(df, year):
    """
    Take a df, slice by year, and produce a list of months,
    where each month is a 2D array in the shape of the calendar
    :param df: dataframe or series
    :return: matrix for daily values and numerals
    """
    df = df[df.index.year == year]

    # Empty matrices
    a = np.empty((6, 7))
    a[:] = np.nan

    day_nums = {m: np.copy(a) for m in range(1, 13)}  # matrix for day numbers
    day_vals = {m: np.copy(a) for m in range(1, 13)}  # matrix for day values

    # Logic to shape datetimes to matrices in calendar layout
    for d in df.iteritems():  # use iterrows if you have a DataFrame

        day = d[0].day
        month = d[0].month
        col = d[0].dayofweek
        first_day_of_the_week = d[0].replace(day=1).weekday()

        row = int((day - 1 + first_day_of_the_week) / 7)

        day_nums[month][row, col] = day  # day number (0-31)
        day_vals[month][row, col] = d[1]  # day value (the heatmap data)

    return day_nums, day_vals


def create_year_calendar( file_name, title, day_nums, day_vals):
    fig, ax = plt.subplots(3, 4, figsize=(14.85, 10.5))

    for i, axs in enumerate(ax.flat):

        axs.imshow(day_vals[i + 1], cmap='viridis', vmin=1, vmax=365)  # heatmap
        axs.set_title(month_names[i])

        # Labels
        axs.set_xticks(np.arange(len(days)))
        axs.set_xticklabels(days, fontsize=10, fontweight='bold', color='#555555')
        axs.set_yticklabels([])

        # Tick marks
        axs.tick_params(axis=u'both', which=u'both', length=0)  # remove tick marks
        axs.xaxis.tick_top()

        # Modify tick locations for proper grid placement
        axs.set_xticks(np.arange(-.5, 6, 1), minor=True)
        axs.set_yticks(np.arange(-.5, 5, 1), minor=True)
        axs.grid(which='minor', color='w', linestyle='-', linewidth=2.1)

        # Despine
        for edge in ['left', 'right', 'bottom', 'top']:
            axs.spines[edge].set_color('#FFFFFF')

        # Annotate
        for w in range(len(weeks)):
            for d in range(len(days)):
                day_val = day_vals[i + 1][w, d]
                day_num = day_nums[i + 1][w, d]

                # Value label
                axs.text(d, w + 0.3, f"{day_val:0.0f}",
                         ha="center", va="center",
                         fontsize=7, color="w", alpha=0.8)

                # If value is 0, draw a grey patch
                if day_val == 0:
                    patch_coords = ((d - 0.5, w - 0.5),
                                    (d - 0.5, w + 0.5),
                                    (d + 0.5, w + 0.5),
                                    (d + 0.5, w - 0.5))

                    square = Polygon(patch_coords, fc='#DDDDDD')
                    axs.add_artist(square)

                # If day number is a valid calendar day, add an annotation
                if not np.isnan(day_num):
                    axs.text(d + 0.45, w - 0.31, f"{day_num:0.0f}",
                             ha="right", va="center",
                             fontsize=6, color="#003333", alpha=0.8)  # day

                # Aesthetic background for calendar day number
                patch_coords = ((d - 0.1, w - 0.5),
                                (d + 0.5, w - 0.5),
                                (d + 0.5, w + 0.1))

                triangle = Polygon(patch_coords, fc='w', alpha=0.7)
                axs.add_artist(triangle)

    # Final adjustments
    fig.suptitle(title, fontsize=16)
    plt.subplots_adjust(left=0.04, right=0.96, top=0.88, bottom=0.04)

    # Save to file
    plt.savefig(file_name)


def create_calendar(name, calendar_title, data):
    """
    Creates a calendar graphic where the user can see the lowest price for east day
    :param file_name: Name to save the calendar to
    :param calendar_title: Name of the calendar
    :param data:
    :return:
    """
    idx, values = zip(*data)
    df = pd.Series(values, idx)
    for year in years:
        day_nums, day_vals = split_months(df, year)
        create_year_calendar(f'{name}-{year}.pdf', calendar_title, day_nums, day_vals)

