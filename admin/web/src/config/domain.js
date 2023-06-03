/**
 * 不同模式下的域名配置
 */

const ENV_CONFIG_MAP = {
	development: 'http://localhost:8089',
	test: 'http://localhost:8089',
	production: 'http://localhost:8089',
}

const domain = ENV_CONFIG_MAP[process.env.NODE_ENV]

export default domain
