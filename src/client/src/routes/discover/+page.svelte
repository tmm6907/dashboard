<script>
    import { onMount } from "svelte";
    import Section from "../../components/Section.svelte";

    let items = $state([]);

    const getDiscoveryItems = async () => {
        try {
            let response = await fetch(
                `https://api.mashboard.app/api/feeds/items/discover/`,
                {
                    credentials: "include",
                },
            );
            if (response.status != 200) {
                throw new Error(await response.text());
            }
            ({ items } = await response.json());
            console.log("discover items:", items);
        } catch (e) {
            console.error(e);
        }
    };

    onMount(async () => {
        console.log("loaded");
        await getDiscoveryItems();
    });
</script>

<div class="p-4">
    <Section heading="Discover" {items} vertical={true} />
</div>
