import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import Cookies from 'js-cookie'
import './MyBooking.css'  // Import the CSS file

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

    const fetchBookings = async () => {
      try {
        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/bookings/my`,
            {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${Cookies.get('token')}`,
                    'Content-Type': 'application/json'
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
							</tr>
						))}
					</tbody>
				</table>
			</div>
		)
  }

export default MyBooking