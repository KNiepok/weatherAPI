// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import weatherapi "github.com/kniepok/weatherAPI"

// WeatherStorage is an autogenerated mock type for the WeatherStorage type
type WeatherStorage struct {
	mock.Mock
}

// GetStatistics provides a mock function with given fields: _a0
func (_m *WeatherStorage) GetStatistics(_a0 *weatherapi.Location) (*weatherapi.Statistics, error) {
	ret := _m.Called(_a0)

	var r0 *weatherapi.Statistics
	if rf, ok := ret.Get(0).(func(*weatherapi.Location) *weatherapi.Statistics); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*weatherapi.Statistics)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*weatherapi.Location) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreWeather provides a mock function with given fields: _a0
func (_m *WeatherStorage) StoreWeather(_a0 *weatherapi.Weather) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*weatherapi.Weather) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
