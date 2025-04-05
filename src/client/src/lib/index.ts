// place files you want to import through the `$lib` alias in this folder.

import { alertState } from "./state.svelte";

export interface FeedData {
    latest: any[]; // Replace 'any' with your specific type if possible
    items: any[];  // Replace 'any' with your specific type if possible
    saved: any[];  // Replace 'any' with your specific type if possible
    collections: any[]; // Replace 'any' with your specific type if possible
}

export const fetchFeedItems = async (category): Promise<FeedData | undefined> => {
    console.log("Fetching feed items");
    category = category ? category : ""
    try {
        const response = await fetch("https://api.mashboard.app/api/feeds/items?category=" + category, {
            credentials: "include",
        });

        if (response.status == 302) {
            console.log("redirect")
            window.location.href = await response.text()
            return
        }
        console.log(response.status)

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const { latest, items, saved, collections } = await response.json() as FeedData;
        return { latest, items, saved, collections };
    } catch (err) {
        console.error("Error fetching feeds:", err);
        return undefined; // Return undefined to indicate failure
    }
};

export const sortFeedItems = (items: any[]) => {
    // sort items on item.pubDate

    return items.sort(
        (a, b) => new Date(b.pub_date) - new Date(a.pub_date)
    );
};

export const getTimeAgoAll = (time: string) => {
    let date = new Date(time);
    const now = new Date();
    const diff = Math.floor((now.getTime() - date.getTime()) / 1000); // Difference in seconds
    const days = Math.floor(diff / (60 * 60 * 24));
    const hours = Math.floor(diff / (60 * 60));
    const minutes = Math.floor(diff / 60);

    if (days > 0) return `${days}d ago`;
    if (hours > 0) return `${hours}h ago`;
    return `${minutes}m ago`;
};

export const getTimeAgo = (time: string) => {
    let date = new Date(time.replace(" ", "T") + "Z");

    const now = new Date();
    const diff = Math.floor((now.getTime() - date.getTime()) / 1000); // Difference in seconds
    const days = Math.floor(diff / (60 * 60 * 24));
    const hours = Math.floor(diff / (60 * 60));
    const minutes = Math.floor(diff / 60);

    if (days > 0) return `${days}d ago`;
    if (hours > 0) return `${hours}h ago`;
    return `${minutes}m ago`;
};

export interface AlertConfig {
    type: string;
    duration: number;
    closable: boolean;
}


export const triggerAlert = (msg: string, config: AlertConfig) => {
    if (msg === "") {
        throw new Error("msg may not be empty")
    }
    alertState.msg = msg;
    if (config.type) {
        alertState.type = config.type
    }
    if (config.duration) {
        alertState.duration = config.duration
    }
    if (config.closable) {
        alertState.closable = config.closable
    }
    alertState.showAlert = true
    setTimeout(() => {
        console.log("closing");
        alertState.showAlert = false;
    }, alertState.duration);
}