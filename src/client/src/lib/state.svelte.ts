import { fetchFeedItems, sortFeedItems } from "$lib";

export const feedState = $state({
    feedItems: [],
    feedCollections: [],
    feedLatest: [],
    feedSaved: [],
    category: "all"
});

export const selectedItem = $state({
    id: 1,
    data: {}
})



