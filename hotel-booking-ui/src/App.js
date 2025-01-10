import React from 'react'
import { BrowserRouter as Router, Route, Routes } from 'react-router'
import ResetPassword from './pages/ResetPassword'
import Login from './pages/Login'
import Register from './pages/Register'
import Home from './pages/Home'
import AdminUsers from './pages/AdminUsers'
import UpdateUserRole from './pages/UpdateUserRole'
import ForgotPassword from './pages/ForgotPassword'
import HotelsAndRooms from './pages/HotelsAndRooms'
import CreateOfflineBooking from './pages/CreateOfflineBooking'
import OwnerHotels from './pages/OwnerHotels'
import CreateHotel from './pages/CreateHotel'
import OwnerRooms from './pages/OwnerRooms'
import CreateRoom from './pages/CreateRoom'
import EditRoom from './pages/EditRoom'
import MyBooking from './pages/MyBooking'


const App = () => {
	return (
		<Router>
			<Routes>
				<Route path='/' element={<Home />} />
				<Route path='/auth/register' element={<Register />} />
				<Route path='/auth/login' element={<Login />} />
				<Route path='/auth/reset-password' element={<ResetPassword />} />
				<Route path='/auth/forgot-password' element={<ForgotPassword />} />
				<Route path='/admin/users' element={<AdminUsers />} />
				<Route path='/admin/users/:id/role' element={<UpdateUserRole />} />
				<Route path='/hotels-and-rooms' element={<HotelsAndRooms />} />
				<Route
					path='/bookings/offline/create'
					element={<CreateOfflineBooking />}
				/>
				<Route path='/owner/hotels' element={<OwnerHotels />} />
        <Route path='/owner/hotels/create' element={<CreateHotel />} />
        <Route path='/owner/rooms' element={<OwnerRooms />} />
        <Route path='/owner/rooms/create' element={<CreateRoom />} />
        <Route path='/owner/rooms/:id/edit' element={<EditRoom />} />
				<Route path='/my-bookings' element={<MyBooking />} />
			</Routes>
		</Router>
	)
}


export default App
