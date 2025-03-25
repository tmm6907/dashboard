// place files you want to import through the `$lib` alias in this folder.


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
        const response = await fetch("http://localhost:8080/api/feeds/items?category=" + category, {
            credentials: "include",
        });

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
        (a, b) => new Date(b.item.pubDate) - new Date(a.item.pubDate),
    );
};