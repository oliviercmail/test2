<template>
  <b-container
    fluid
    class="d-flex justify-content-center vh-100 logo"
  >
    MAP
  </b-container>
</template>

<script>
export default {
  props: {
    logo: {
      type: String,
      default: () => ''
    },
  },
}
</script>

<style lang="scss" scoped>
.logo {
  img {
    animation: pulse 4.2s infinite;
  }

  @keyframes pulse {
    0% {
      opacity: 0;
    }
    65% {
      opacity: 1;
    }
    100% {
      opacity: 0;
    }
  }
}
</style>
