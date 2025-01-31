<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import AnimatedGradient from "$lib/components/AnimatedGradient.svelte";
</script>
<div class="container flex flex-col items-center justify-center w-full h-screen max-w-none px-0">
    <AnimatedGradient 
        duration="8s"
        colors={{
            from: "#8B5CF6",
            via: "#EC4899",
            to: "#EF4444"
        }}
        degree="-45deg"
    > 
        <div class="flex flex-col items-center justify-center">
            <h1 class="text-4xl">Janko Muzykant</h1>
            <div class="mt-10"></div>
            <div class="flex gap-2 w-full">
                <a href="/auth" class="w-full">
                    <Button class="w-full">Login</Button>
                </a>
                <a href="/signup" class="w-full">
                    <Button class="w-full">Signup</Button>
                </a>
            </div>
        </div>
    </AnimatedGradient>
</div>