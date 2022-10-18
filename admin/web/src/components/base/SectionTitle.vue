<template>
  <div class="section-title">
    <span class="section-title__tag"></span>
    <span class="section-title__name">{{ name }}</span>
  </div>
</template>

<script>
export default {
  props: {
    name: {
      type: String,
      default: ''
    }
  }
}
</script>

<style lang="scss" scoped>
.section-title {
  display: inline-block;
  margin-right: 10px;

  .section-title__tag {
    display: inline-block;
    width: 6px;
    height: 18px;
    background-color: var(--theme);
    border-radius: 50%;
    vertical-align: middle;
  }

  .section-title__name {
    font-size: 18px;
    margin-left: 10px;
    vertical-align: middle;
  }
}
</style>
