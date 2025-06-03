import { MiddlewareConsumer, Module } from '@nestjs/common';
import { LtiService } from './lti.service';
import { LtiController } from './lti.controller';
import { LtiMiddleware } from './lti.middleware';

@Module({
  controllers: [LtiController],
  providers: [LtiService],
})
export class LtiModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(LtiMiddleware).forRoutes('/lti');
  }
}
