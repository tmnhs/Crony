<template>
  <!-- 富文本编辑器 -->
  <textarea
    :id="id"
    style="visibility:hidden"
    v-loading="!hasInit">
   </textarea>
</template>

<script>
  import { loadScript } from '@/utils/core';

  const tinymceCDN = 'https://cdn.jsdelivr.net/npm/tinymce@5.1.2/tinymce.min.js';

  export default {
    props: {
      // 编辑器id
      id: {
        type: String,
        default () {
          // 如果一同个页面中要使用多个编辑器，id需要不一样。
          return 'tinymce-' + +new Date() + ((Math.random() * 1000).toFixed(0) + '');
        }
      },
      // 内容
      value: {
        type: String,
        default: ''
      },
      // 编辑器高度
      height: {
        type: Number,
        default: 500
      }
    },
    data() {
      return {
        hasInit: false,
        hasInput: false
      };
    },
    watch: {
      // 进行编辑的时候不需要手动设置新的内容.
      // 在编辑器初始化完成后才能设置。
      // 为什么使用insertContent而不使用setContent？setContent之后光标会跑到内容最前方。
      value(newVal) {
        if (!this.hasInput && this.hasInit) {
          tinymce.get(this.id).insertContent(newVal);
        }
      }
    },
    mounted() {
      this.init();
    },
    // 组件销毁时必须手动销毁tinymce。比如打开了工具栏，然后关闭当前页面或者跳转到其他页面，由于编辑器没有被销毁，这个工具栏就会一直悬浮在页面上。
    destroyed() {
      this.destoryTinymce();
    },
    // 如果使用tinymce的页面没有被keep-alive缓存，则不需要activated和deactivated这两个钩子。
    // 在缓存的页面间切换时必须先销毁编辑器实例再重新初始化。
    activated() {
      this.init();
    },
    deactivated() {
      this.destoryTinymce();
    },
    methods: {
      init() {
        loadScript(tinymceCDN, (err, res) => {
          if (!err) {
            this.initTinymce();
          } else {
            this.$message.error(err.message);
          }
        })
      },
      initTinymce() {
        tinymce.init({
          // 挂载元素
          selector: `#${this.id}`,
          // 语言
          language: 'zh_CN',
          // 指定语言文件位置
          language_url: '/static/tinymce/langs/zh_CN.js',
          // 编辑器高度
          height: this.height,
          // 要使用的插件
          plugins: 'advlist autosave autolink charmap code codesample emoticons fullscreen help hr image insertdatetime link lists media pagebreak paste preview print searchreplace table textpattern visualblocks visualchars wordcount',
          // 设置在工具栏中显示的操作按钮
          toolbar: 'restoredraft | undo redo | styleselect | bold italic underline forecolor backcolor removeformat | alignleft aligncenter alignright alignjustify | bullist numlist outdent indent | table link image media emoticons charmap | searchreplace fullscreen preview code print',
          // 右键菜单
          contextmenu: 'copy paste link inserttable ',
          // 自动获得焦点
          auto_focus: this.id,
          // 带格式粘贴,比如从word文档中复制了一段内容来粘贴,则不会清除格式.
          paste_enable_default_filters: false,
          // 定义插入的时间格式
          insertdatetime_formats: ['%Y-%m-%d %H:%M:%S', '%Y/%m/%d %H:%M:%S', '%Y年%m月%d日 %H点%M分%S秒', '%Y年%m月%d日', '%H点%M分'],
          // 表格操作工具栏
          table_toolbar: "tableprops tablerowprops tablecellprops | tableinsertcolbefore tableinsertcolafter tabledeletecol | tableinsertrowbefore tableinsertrowafter tabledeleterow | tablemergecells tablesplitcells | tabledelete",
          // 点击超链接时显示快捷工具条。用户可能不知道要按住ctrl再点击来打开超链接
          link_context_toolbar: true,
          // 为超链接提示是否添加http://前缀。如果是fasle，只会为前缀为www的超链接进行提示。
          link_assume_external_targets: true,
          // 编辑器初始化之前执行的钩子
          setup: editor => {},
          // 编辑器初始化完毕的钩子。如果不使用箭头函数，内部的this指向的是当前editor实例。
          init_instance_callback: editor => {
            this.hasInit = true;
            if (this.value) {
              editor.insertContent(this.value);
            };
            editor.on('setContent keyup', () => {
              this.hasInput = true;
              this.$emit('input', editor.getContent());
            });
          },
          // 插入图片时的钩子
          images_upload_handler: (blobInfo, success, failure) => {
            // 图片file对象
            const file = blobInfo.blob();
            // 图片的base64数据
            let base64 = blobInfo.base64();
            const fileType = file.type;
            if (fileType === 'image/jpeg') {
              base64 = 'data:image/jpeg;base64,' + base64;
            } else if (fileType === 'image/png') {
              base64 = 'data:image/png;base64,' + base64;
            } else if (fileType === 'image/gif') {
              base64 = 'data:image/gif;base64,' + base64;
            }
            const formData = new FormData();
            formData.append('file', file);
            formData.append('name', file.name);
            // 如果后端需要在文件服务中保存这个图片，就将图片传给后端
            // this.$axios({
            //     url: '',
            //     method: 'post',
            //     data: formData
            // }).then(res => {
            //     const imgURL = res.url;
            //     success(imgURL);
            // })
            // success的参数需要是远程图片地址，传本地的图片进去如果需要对图片进行编辑，比如拖拉进行缩放会报错图片不存在，因为使用的是blobUri()方法得到的地址。这里只是展示插入图片后的效果
            success(base64);
          }
        });
      },
      destoryTinymce() {
        // 关闭当前页面时，deactivated钩子会在destroyed之前执行，此时editor实例已经被销毁了。
        if (tinymce.get(this.id)) {
          tinymce.get(this.id).destroy();
        }
      }
    }
  };
</script>
