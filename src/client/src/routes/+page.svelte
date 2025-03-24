<script>
    import { onMount } from "svelte";
    import BottomNav from "../components/BottomNav.svelte";
    import Feeds from "../components/FeedItems.svelte";

    import Filter from "../components/Filter.svelte";
    import SubscribeForm from "../components/SubscribeForm.svelte";
    import FeedItems from "../components/FeedItems.svelte";

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

<Filter />
<FeedItems />
<SubscribeForm />
<BottomNav />
