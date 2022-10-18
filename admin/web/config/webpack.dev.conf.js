const path = require('path')
const { merge } = require('webpack-merge')
const baseWebpackConfig = require('./webpack.base.conf.js')
const webpack = require('webpack')

const resolve = dir => {
	return path.resolve(process.cwd(), dir)
}

module.exports = merge(baseWebpackConfig, {
	mode: 'development',
	plugins: [new webpack.HotModuleReplacementPlugin()],
	devServer: {
		contentBase: resolve('dist'),
		compress: true,
		port:5000,
		open: true,
		hot: true,
		overlay: true,
	},
	module: {
		rules: [
			{
				test: /\.css$/,
				use: ['style-loader', 'css-loader'],
			},
			{
				test: /\.scss$/,
				use: [
					'style-loader',
					'css-loader',
					{
						loader: 'sass-loader',
						options: {
							data: `@import "./src/assets/styles/variable.scss";@import "./src/assets/styles/mixin.scss";`,
						},
					},
				],
			},
		],
	},
})
