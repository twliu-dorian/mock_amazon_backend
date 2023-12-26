package user

type Controller struct{}

func (Controller) Create(c echo.Context) (err error) {
	var request CreateUserRequest
	if err = c.Bind(&request); err != nil {
		err = new(apierror.ApiError).From(err).SetCode(apierror.INPUT_ERROR)
		return
	}
	if err = request.Validate(); err != nil {
		err = new(apierror.ApiError).From(err).SetCode(apierror.INPUT_ERROR)
		return
	}

	toCreate := request.ToWallet()
	err = CreateUser(toCreate)
	if err != nil {
		err = new(apierror.ApiError).From(err)
		return
	}
	
	return util.AcceptedResponse(c, taskId)
}