<template>
  <div class="menu-item">
    <el-submenu v-if="config.children" :index="config.path">
      <template slot="title">
        <base-icon v-if="config.icon" :name="config.icon" icon-class="menu-icon"></base-icon>
        <span>{{ config.title }}</span>
      </template>
      <menu-item
        v-for="childConfig in config.children"
        :key="childConfig.path"
        :config="childConfig"
      ></menu-item>
    </el-submenu>

    <el-menu-item v-else :index="config.path">
      <base-icon v-if="config.icon" :name="config.icon" icon-class="menu-icon"></base-icon>
      <span slot="title">{{ config.title }}</span>
    </el-menu-item>
  </div>
</template>

<script>
export default {
  name: 'MenuItem',
  props: {
    config: {
      required: true,
      default: {}
    }
  }
}
</script>

<style lang="scss" scoped>
.menu-item {
  .menu-icon {
    vertical-align: middle;
    font-size: 18px;
  }
}
</style>

<style lang="scss">
.el-menu--collapse {
  .menu-item {
    .el-submenu > .el-submenu__title {
      span,
      .el-submenu__icon-arrow {
        display: none;
      }
    }
  }
}
</style>
