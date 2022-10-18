export default {
	inserted(el) {
		const drop = el.querySelector('.el-dialog__header')
		const dialog = el.querySelector('.el-dialog')
		drop.style.cursor = 'move'
		drop.onmousedown = event => {
			//鼠标按下时，鼠标到弹框外边框的距离，这个距离是固定的
			const x = event.clientX - dialog.offsetLeft
			const y = event.clientY - dialog.offsetTop

			// 弹框的左外边距和上外边距
			// const dialogML = (document.body.clientWidth - dialog.offsetWidth) / 2;
			const dialogML = parseInt(window.getComputedStyle(dialog).marginLeft)
			const dialogMT = parseInt(window.getComputedStyle(dialog).marginTop)

			// 当前视窗的大小
			const screenWidth = window.innerWidth
			const screenHeight = window.innerHeight

			const minDialogLeft = -dialogML
			const maxDialogLeft = dialogML
			const minDialogTop = -dialogMT
			const maxDialogTop = screenHeight - dialog.offsetHeight - dialogMT
			// console.log(maxDialogTop);

			document.onmousemove = event => {
				let styleLeft = event.clientX - x - dialogML
				let styleTop = event.clientY - y - dialogMT
				// console.log(styleTop);

				if (styleLeft <= minDialogLeft) {
					styleLeft = minDialogLeft
				} else if (styleLeft >= maxDialogLeft) {
					styleLeft = maxDialogLeft
				}
				if (styleTop <= minDialogTop) {
					styleTop = minDialogTop
				} else if (styleTop >= maxDialogTop) {
					styleTop = maxDialogTop
				}

				dialog.style.left = styleLeft + 'px'
				dialog.style.top = styleTop + 'px'
				window.getSelection ? window.getSelection().removeAllRanges() : document.selection.empty() //清除选中的文本（点击带header中的文字进行拖动可能造成选中文本）
			}
		}
		document.onmouseup = function () {
			document.onmousemove = null
		}
	},
}
