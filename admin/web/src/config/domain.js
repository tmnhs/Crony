/**
 * 不同模式下的域名配置
 */

const ENV_CONFIG_MAP = {
	development: 'http://localhost:8088',
	test: 'http://localhost:8088',
	production: 'http://114.55.178.217:8089',
}

const domain = ENV_CONFIG_MAP[process.env.NODE_ENV]

export default domain
