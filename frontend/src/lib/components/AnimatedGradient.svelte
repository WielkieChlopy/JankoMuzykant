<script lang="ts">
    let { children, duration, colors, degree } = $props();
    // let duration = "8s";
    // let colors = {
    //   from: "#8B5CF6",  // purple-500
    //   via: "#EC4899",   // pink-500
    //   to: "#EF4444"     // red-500
    // };
    // let degree = "-45deg";
  </script>
  
  <div class="gradient-container" style="--color-from: {colors.from}; --color-via: {colors.via}; --color-to: {colors.to}; --animation-duration: {duration}; --degree: {degree};">
    <div class="gradient-content">
      {@render children()}
    </div>
  </div>
  
  <style>
    .gradient-container {
      position: relative;
      height: 100%;  
      width: 100%;
      overflow: hidden;
      background: linear-gradient(
        var(--degree),
        var(--from-color),
        var(--via-color),
        var(--to-color)
      );
      background-size: 200% 200%;
      animation: gradient var(--animation-duration) ease infinite;
    }
  
    .gradient-content {
      position: absolute;
      inset: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
    }
  
    @keyframes gradient {
      0% {
        background-position: 0% 50%;
      }
      50% {
        background-position: 100% 50%;
      }
      100% {
        background-position: 0% 50%;
      }
    }
  
    :root {
      --from-color: var(--color-from, #8B5CF6);
      --via-color: var(--color-via, #EC4899);
      --to-color: var(--color-to, #EF4444);
      --animation-duration: var(--duration, 8s);
      --degree: var(--degree, -45deg);
    }
  </style>