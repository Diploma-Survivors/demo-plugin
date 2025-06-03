import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { LtiModule } from './lti/lti.module';
import { ConfigModule } from '@nestjs/config';
import * as Joi from 'joi';

@Module({
  imports: [
    LtiModule,
    ConfigModule.forRoot({
      validationSchema: Joi.object({
        DATABASE_NAME: Joi.string().required(),
        DATABASE_USERNAME: Joi.string().required(),
        DATABASE_PASSWORD: Joi.string().required(),
        DATABASE_PORT: Joi.string().required(),
        DATABASE_URI: Joi.string().required(),
        PORT: Joi.string().required(),
        LTI_KEY: Joi.string().required(),
        LTI_NAME: Joi.string().required(),
        LTI_PLATFORM_URL: Joi.string().required(),
        LTI_CLIENT_ID: Joi.string().required(),
        LTI_PUBLIC_KEYSET_URL: Joi.string().required(),
        LTI_ACCESS_TOKEN_URL: Joi.string().required(),
        LTI_AUTHENTICATION_URL: Joi.string().required(),
      }),
      isGlobal: true,
    }),
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
