import { Controller, Get, Req, Res } from '@nestjs/common';
import { LtiService } from './lti.service';
import { Request, Response } from 'express';

@Controller('lti')
export class LtiController {
  constructor(private readonly ltiService: LtiService) {}

  @Get('nolti')
  nolti(@Req() req: Request, @Res() res: Response) {
    res.send(
      'There was a problem getting you authenticated with the attendance application. Please contact support.',
    );
  }

  @Get('ping')
  ping(@Req() req: Request, @Res() res: Response) {
    res.send('pong');
  }

  @Get('protected')
  protected(@Req() req: Request, @Res() res: Response) {
    res.send('Insecure');
  }
}
