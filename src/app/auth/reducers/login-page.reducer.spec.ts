import { reducer } from '../reducers/login-page.reducer';
import * as fromLoginPage from '../reducers/login-page.reducer';

import { AuthApiActions, LoginPageActions } from '../actions';

import { Credentials } from '../models/credentials';
import { Account } from '../../lib/model';

describe('LoginPageReducer', () => {
  describe('undefined action', () => {
    it('should return the default state', () => {
      const action = {} as any;

      const result = reducer(undefined, action);

      expect(result).toMatchSnapshot();
    });
  });

  describe('LOGIN', () => {
    it('should make pending to true', () => {
      const credentials = { passphrase: 'test' } as Credentials;
      const createAction = new LoginPageActions.Login({ credentials });

      const expectedResult = {
        error: null,
        pending: true,
      };

      const result = reducer(fromLoginPage.initialState, createAction);

      expect(result).toMatchSnapshot();
    });
  });

  describe('LOGIN_SUCCESS', () => {
    it('should have no error and no pending state', () => {
      const account = new Account();
      const createAction = new AuthApiActions.LoginSuccess({ account });

      const expectedResult = {
        error: null,
        pending: false,
      };

      const result = reducer(fromLoginPage.initialState, createAction);

      expect(result).toMatchSnapshot();
    });
  });

  describe('LOGIN_FAILURE', () => {
    it('should have an error and no pending state', () => {
      const error = 'login failed';
      const createAction = new AuthApiActions.LoginFailure({ error });

      const expectedResult = {
        error: error,
        pending: false,
      };

      const result = reducer(fromLoginPage.initialState, createAction);

      expect(result).toMatchSnapshot();
    });
  });
});
