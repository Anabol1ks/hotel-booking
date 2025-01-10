import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import Cookies from 'js-cookie'
import './MyBooking.css'

const MyBooking = () => {
	const [bookings, setBookings] = useState([])
	const [error, setError] = useState('')
	const navigate = useNavigate()

	useEffect(() => {
		const token = Cookies.get('token')
		if (token) {
			fetchBookings()
		} else {
			navigate('auth/login')
		}
	}, [navigate])

	const calculateRemainingTime = createdAt => {
		const bookingTime = new Date(createdAt).getTime()
		const currentTime = new Date().getTime()
		const timeDiff = bookingTime + 30 * 60 * 1000 - currentTime // 30 minutes

		if (timeDiff <= 0) {
			return null // Time expired
		}

		const minutes = Math.floor((timeDiff / (1000 * 60)) % 60)
		const seconds = Math.floor((timeDiff / 1000) % 60)

		return `${minutes}:${seconds < 10 ? '0' : ''}${seconds}`
	}

	const fetchBookings = async () => {
		try {
			const response = await fetch(
				`${process.env.REACT_APP_API_URL}/bookings/my`,
				{
					method: 'GET',
					headers: {
						Authorization: `Bearer ${Cookies.get('token')}`,
						'Content-Type': 'application/json',
					},
				}
			)

			if (response.ok) {
				const data = await response.json()
				setBookings(data)
			} else {
				setError('Не удалось загрузить бронирования')
			}
		} catch (err) {
			setError('Ошибка при загрузке бронирований')
		}
	}

	const handlePayBooking = async bookingId => {
		try {
			const response = await fetch(
				`${process.env.REACT_APP_API_URL}/bookings/${bookingId}/pay`,
				{
					method: 'POST',
					headers: {
						Authorization: `Bearer ${Cookies.get('token')}`,
						'Content-Type': 'application/json',
					},
				}
			)

			if (response.ok) {
				const data = await response.json()
				window.location.href = data.payment_url
			} else {
				const errorData = await response.json()
				setError(errorData.error || 'Не удалось создать платеж')
			}
		} catch (err) {
			setError('Ошибка при создании платежа')
		}
	}

	const getPaymentStatusLabel = status => {
		switch (status) {
			case 'succeeded':
				return 'Оплачен'
			case 'pending':
				return 'Не оплачен'
			default:
				return 'Статус не определен'
		}
	}

	const PendingBookingTimer = ({ booking }) => {
		const [remainingTime, setRemainingTime] = useState(
			calculateRemainingTime(booking.CreatedAt)
		)

		useEffect(() => {
			const timer = setInterval(() => {
				const newRemainingTime = calculateRemainingTime(booking.CreatedAt)
				setRemainingTime(newRemainingTime)
			}, 1000)

			return () => clearInterval(timer)
		}, [booking.CreatedAt])

		if (!remainingTime) {
			return <span className='expired-timer'>Время оплаты истекло</span>
		}

		return <span className='active-timer'>{remainingTime}</span>
	}

	return (
		<div className='my-bookings-container'>
			<h1 className='my-bookings-title'>Мои бронирования</h1>
			{error && <p className='error-message'>{error}</p>}
			<table className='bookings-table'>
				<thead>
					<tr>
						<th>ID</th>
						<th>Номер комнаты</th>
						<th>Дата заезда</th>
						<th>Дата выезда</th>
						<th>Общая стоимость</th>
						<th>Статус оплаты</th>
						<th>Время до оплаты</th>
						<th>Действия</th>
					</tr>
				</thead>
				<tbody>
					{bookings.map(booking => (
						<tr key={booking.ID}>
							<td>{booking.ID}</td>
							<td>{booking.RoomID}</td>
							<td>{new Date(booking.StartDate).toLocaleDateString()}</td>
							<td>{new Date(booking.EndDate).toLocaleDateString()}</td>
							<td>{booking.TotalCost} руб.</td>
							<td>{getPaymentStatusLabel(booking.PaymentStatus)}</td>
							<td>
								{booking.PaymentStatus === 'pending' && (
									<PendingBookingTimer booking={booking} />
								)}
							</td>
							<td>
								{booking.PaymentStatus === 'pending' && (
									<button
										onClick={() => handlePayBooking(booking.ID)}
										className='pay-button'
									>
										Оплатить
									</button>
								)}
							</td>
						</tr>
					))}
				</tbody>
			</table>
		</div>
	)
}

export default MyBooking
