import { LtiMiddleware } from './lti.middleware';

describe('LtiMiddleware', () => {
  it('should be defined', () => {
    expect(new LtiMiddleware()).toBeDefined();
  });
});
