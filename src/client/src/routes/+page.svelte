<script>
    import { onMount } from "svelte";
    import BottomNav from "../components/BottomNav.svelte";
    import SubscribeForm from "../components/SubscribeForm.svelte";
    import Feed from "../components/Feed.svelte";

    import { fetchFeedItems, sortFeedItems } from "$lib";
    import { feedState } from "$lib/state.svelte";
    $effect(async () => {
        let results = await fetchFeedItems(feedState.category);
        console.log("category changed", feedState.category);
        feedState.feedItems = sortFeedItems(results.items ? results.items : []);
        feedState.feedCollections = sortFeedItems(
            results.collections ? results.collections : [],
        );
        feedState.feedLatest = sortFeedItems(
            results.latest ? results.latest : [],
        );
        feedState.feedSaved = sortFeedItems(results.saved ? results.saved : []);
    });

    onMount(() => {
        const login = async () => {
            try {
                const response = await fetch("http://localhost:8080/auth/", {
                    method: "GET",
                    credentials: "include", // Important for cookies/sessions
                });

                if (response.status === 302) {
                    const redirectUrl = await response.text();
                    if (redirectUrl) {
                        window.location.href = redirectUrl; // Redirect manually
                    }
                }
            } catch (e) {
                console.error(e);
            }
        };

        login();
    });
</script>

<Feed />
<SubscribeForm />
<BottomNav />
