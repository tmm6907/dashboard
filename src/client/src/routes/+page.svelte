<script>
    import { fetchFeedItems, sortFeedItems, triggerAlert } from "$lib";
    import { feedState } from "$lib/state.svelte";
    import { onMount } from "svelte";
    import Feed from "../components/Feed.svelte";
    import Filter from "../components/Filter.svelte";
    import { PUBLIC_API_BASE_URL, PUBLIC_HOSTNAME } from "$env/static/public";
    async function updateFeed() {
        let results = await fetchFeedItems(feedState.category);
        console.log("category changed", feedState.category);
        console.log(results);
        if (!results) {
            console.error("failed to update feed");
            return;
        }
        feedState.feedItems = sortFeedItems(results.items);
        feedState.feedCollections = results.collections;
        feedState.feedLatest = sortFeedItems(results.latest);
    }

    const authenticateUser = async () => {
        let resp = await fetch("https://api.mashboard.app/auth/login");
        if (!resp) {
            console.error("failed to authenticate");
            throw Error("failed to authenticate");
        }
        if (!resp.ok) {
            console.error("failed to authenticate:", resp.status);
            throw new Error("failed to authenticate");
        }
    };

    $effect(async () => {
        try {
            await authenticateUser();
        } catch (e) {
            console.error(e);
            triggerAlert({ e, type: "alert-error" });
        }
        await updateFeed();
    });
    onMount(async () => {
        await updateFeed();
    });
</script>

<Filter />
<Feed />
