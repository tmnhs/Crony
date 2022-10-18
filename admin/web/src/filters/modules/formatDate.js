import dayjs from 'dayjs'

export default function formatDate(date, format = 'YYYY-MM-DD HH:mm:ss') {
	let dateValue = date
	if (!dateValue) {
		dateValue = new Date()
	}
	return dayjs(dateValue).format(format)
}
